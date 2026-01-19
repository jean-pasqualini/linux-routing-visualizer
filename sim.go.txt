// iptsim.go
//
// A single-file, CLI iptables (classic) simulator in Go.
// Focus: chain traversal + jump vs goto semantics, RETURN behavior,
// and a practical subset of matches (proto, src/dst CIDR, sport/dport,
// in/out iface, tcp flags, ctstate, mark, comment).
//
// This is NOT a kernel firewall and not a full iptables parser.
// It’s a deterministic simulator you can use for learning/testing.
//
// Build:
//   go build -o iptsim iptsim.go
//
// Quick start:
//
// 1) Create rules file (example.rules):
// ------------------------------------------------------------
// *filter
// :INPUT DROP [0:0]
// :FORWARD DROP [0:0]
// :OUTPUT ACCEPT [0:0]
// :MYCHAIN - [0:0]
//
// -A INPUT -p tcp -s 10.0.0.0/8 --dport 22 -j ACCEPT
// -A INPUT -p tcp -s 192.168.1.0/24 --dport 22 -j MYCHAIN
// -A MYCHAIN -m conntrack --ctstate NEW -j LOG --log-prefix "ssh new "
// -A MYCHAIN -j RETURN
// -A INPUT -p tcp --dport 80 -g MYCHAIN
// -A INPUT -j DROP
// COMMIT
// ------------------------------------------------------------
//
// 2) List:
//   ./iptsim list -f example.rules
//
// 3) Simulate a packet:
//   ./iptsim sim -f example.rules --table filter --chain INPUT \
//     --proto tcp --src 192.168.1.10 --dst 203.0.113.5 --dport 22 \
//     --in eth0 --ctstate NEW
//
// 4) Explain jump vs goto:
//   In the simulator, -j pushes a return address. -g does not.
//   If a packet hits RETURN inside a chain reached by -g, it returns
//   to the chain that called the *current* chain (i.e., skips the caller
//   chain’s remaining rules). This matches iptables’ documented behavior.
//
// Commands:
//   list  - show parsed tables/chains/rules
//   sim   - simulate a packet through one chain
//
// Notes:
// - Default tables: filter (INPUT/FORWARD/OUTPUT) if not declared.
// - Builtin chain policies honored (ACCEPT/DROP/RETURN unsupported as policy).
// - Targets supported: ACCEPT, DROP, REJECT, RETURN, LOG, MARK, JUMP(userchain), GOTO(userchain).
// - LOG only appends to result logs (no effect on verdict).
// - REJECT treated like DROP but marked as "REJECT".
//
// Enjoy.

package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Verdict string

const (
	VAccept Verdict = "ACCEPT"
	VDrop   Verdict = "DROP"
	VReject Verdict = "REJECT"
	VReturn Verdict = "RETURN" // internal control-flow signal
	VNone   Verdict = "NONE"   // rule didn't decide fate
)

type TableName string

const (
	TableFilter TableName = "filter"
	TableNat    TableName = "nat"
	TableMangle TableName = "mangle"
	TableRaw    TableName = "raw"
	TableSec    TableName = "security"
)

type ChainType int

const (
	ChainBuiltin ChainType = iota
	ChainUser
)

type Chain struct {
	Name     string
	Type     ChainType
	Policy   Verdict // for builtin chains only; empty for user chains (interpreted as RETURN at end)
	Rules    []*Rule
	Counters [2]uint64 // packets, bytes (simulated)
}

type Table struct {
	Name   TableName
	Chains map[string]*Chain
}

type Rule struct {
	LineNo   int
	Raw      string
	Match    Match
	Target   Target
	Counters [2]uint64
}

type Match struct {
	Proto        string // tcp/udp/icmp/any("")
	SrcCIDR      *net.IPNet
	DstCIDR      *net.IPNet
	InIface      string
	OutIface     string
	SrcPort      *PortSpec
	DstPort      *PortSpec
	CTStates     map[string]bool
	MarkMask     uint32
	MarkValue    uint32
	HasMarkMatch bool
	TCPFlags     *TCPFlagsSpec
	Comment      string
}

type TCPFlagsSpec struct {
	Mask map[string]bool
	Comp map[string]bool
}

type PortSpec struct {
	Any    bool
	Single int
	Range  bool
	From   int
	To     int
}

type TargetKind int

const (
	TUnknown TargetKind = iota
	TAccept
	TDrop
	TReject
	TReturn
	TLog
	TMark
	TJump // -j user chain
	TGoto // -g user chain
)

type Target struct {
	Kind        TargetKind
	ChainName   string // for jump/goto
	LogPrefix   string // for LOG
	SetMark     uint32 // for MARK --set-mark
	SetMarkMask uint32
	Raw         string
}

type Packet struct {
	Proto    string
	SrcIP    net.IP
	DstIP    net.IP
	SrcPort  int
	DstPort  int
	InIface  string
	OutIface string
	CTState  string
	Mark     uint32
	TCPFlags map[string]bool
	Bytes    uint64
}

type SimResult struct {
	Table     string   `json:"table"`
	Chain     string   `json:"chain"`
	Verdict   Verdict  `json:"verdict"`
	Reason    string   `json:"reason"`
	Trace     []string `json:"trace"`
	Logs      []string `json:"logs"`
	FinalMark uint32   `json:"final_mark"`
}

type Firewall struct {
	Tables map[TableName]*Table
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "list":
		cmdList(os.Args[2:])
	case "sim":
		cmdSim(os.Args[2:])
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `iptsim - iptables simulator (single-file Go)

Usage:
  iptsim list -f rules.txt
  iptsim sim  -f rules.txt --table filter --chain INPUT [packet flags...]

Packet flags (sim):
  --proto tcp|udp|icmp|any
  --src  IP[/CIDR]   (packet is a single IP; CIDR allowed but will use the IP portion)
  --dst  IP[/CIDR]
  --sport N
  --dport N
  --in  IFACE
  --out IFACE
  --ctstate NEW|ESTABLISHED|RELATED|INVALID
  --mark 0xNN or decimal
  --tcpflags SYN,ACK,FIN,RST,PSH,URG (comma list)
  --bytes N

Other:
  --json    (sim: print JSON result)
`)
}

func cmdList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	file := fs.String("f", "", "rules file")
	fs.Parse(args)

	if *file == "" {
		fmt.Fprintln(os.Stderr, "list: -f is required")
		os.Exit(2)
	}

	fw, err := LoadRules(*file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	printFirewall(fw)
}

func cmdSim(args []string) {
	fs := flag.NewFlagSet("sim", flag.ExitOnError)
	file := fs.String("f", "", "rules file")
	table := fs.String("table", "filter", "table name")
	chain := fs.String("chain", "INPUT", "chain name")

	proto := fs.String("proto", "any", "packet protocol")
	src := fs.String("src", "0.0.0.0", "source IP")
	dst := fs.String("dst", "0.0.0.0", "destination IP")
	sport := fs.Int("sport", 0, "source port")
	dport := fs.Int("dport", 0, "destination port")
	inIf := fs.String("in", "", "input interface")
	outIf := fs.String("out", "", "output interface")
	ctstate := fs.String("ctstate", "", "conntrack state")
	mark := fs.String("mark", "0", "packet mark (hex like 0x1 or decimal)")
	tcpflags := fs.String("tcpflags", "", "tcp flags comma list (SYN,ACK,FIN,RST,PSH,URG)")
	bytes := fs.Uint64("bytes", 60, "packet length in bytes")

	jsonOut := fs.Bool("json", false, "print JSON")

	fs.Parse(args)

	if *file == "" {
		fmt.Fprintln(os.Stderr, "sim: -f is required")
		os.Exit(2)
	}

	fw, err := LoadRules(*file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	srcIP, err := parseIPLoose(*src)
	if err != nil {
		fmt.Fprintln(os.Stderr, "bad --src:", err)
		os.Exit(2)
	}
	dstIP, err := parseIPLoose(*dst)
	if err != nil {
		fmt.Fprintln(os.Stderr, "bad --dst:", err)
		os.Exit(2)
	}

	markVal, err := parseUint32(*mark)
	if err != nil {
		fmt.Fprintln(os.Stderr, "bad --mark:", err)
		os.Exit(2)
	}

	flagsMap := map[string]bool{}
	if *tcpflags != "" {
		for _, f := range strings.Split(*tcpflags, ",") {
			ff := strings.ToUpper(strings.TrimSpace(f))
			if ff == "" {
				continue
			}
			flagsMap[ff] = true
		}
	}

	p := Packet{
		Proto:    strings.ToLower(*proto),
		SrcIP:    srcIP,
		DstIP:    dstIP,
		SrcPort:  *sport,
		DstPort:  *dport,
		InIface:  *inIf,
		OutIface: *outIf,
		CTState:  strings.ToUpper(*ctstate),
		Mark:     markVal,
		TCPFlags: flagsMap,
		Bytes:    *bytes,
	}

	res, err := fw.Simulate(TableName(strings.ToLower(*table)), *chain, p)
	if err != nil {
		fmt.Fprintln(os.Stderr, "simulate error:", err)
		os.Exit(1)
	}

	if *jsonOut {
		b, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(b))
		return
	}

	fmt.Printf("Table=%s Chain=%s Verdict=%s\n", res.Table, res.Chain, res.Verdict)
	if res.Reason != "" {
		fmt.Println("Reason:", res.Reason)
	}
	fmt.Printf("Final mark: 0x%X (%d)\n", res.FinalMark, res.FinalMark)

	if len(res.Logs) > 0 {
		fmt.Println("\nLogs:")
		for _, l := range res.Logs {
			fmt.Println("  " + l)
		}
	}

	fmt.Println("\nTrace:")
	for _, t := range res.Trace {
		fmt.Println("  " + t)
	}
}

// ------------------------------ Loading / Parsing ------------------------------

func LoadRules(path string) (*Firewall, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fw := &Firewall{Tables: map[TableName]*Table{}}
	ensureTable(fw, TableFilter)

	var curTable *Table
	lineNo := 0

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// table header: *filter
		if strings.HasPrefix(line, "*") {
			tname := strings.ToLower(strings.TrimPrefix(line, "*"))
			curTable = ensureTable(fw, TableName(tname))
			continue
		}

		if strings.EqualFold(line, "COMMIT") {
			curTable = nil
			continue
		}

		// chain def: :INPUT DROP [0:0]
		if strings.HasPrefix(line, ":") {
			if curTable == nil {
				curTable = ensureTable(fw, TableFilter)
			}
			if err := parseChainDef(curTable, line); err != nil {
				return nil, fmt.Errorf("line %d: %w", lineNo, err)
			}
			continue
		}

		// rule: -A CHAIN ...
		if strings.HasPrefix(line, "-A ") || strings.HasPrefix(line, "-I ") {
			if curTable == nil {
				curTable = ensureTable(fw, TableFilter)
			}
			r, chainName, err := parseRule(lineNo, line)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", lineNo, err)
			}
			ch := ensureChain(curTable, chainName, ChainUser, "")
			ch.Rules = append(ch.Rules, r)
			continue
		}

		// ignore -P policy lines (iptables-restore uses :CHAIN POLICY)
		if strings.HasPrefix(line, "-P ") {
			if curTable == nil {
				curTable = ensureTable(fw, TableFilter)
			}
			if err := parsePolicyLine(curTable, line); err != nil {
				return nil, fmt.Errorf("line %d: %w", lineNo, err)
			}
			continue
		}

		// ignore other lines (like -F, -X, etc.) for simplicity
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	// Ensure standard builtin chains exist for filter if not provided.
	ensureBuiltinDefaults(fw)

	return fw, nil
}

func ensureTable(fw *Firewall, name TableName) *Table {
	if t, ok := fw.Tables[name]; ok {
		return t
	}
	t := &Table{Name: name, Chains: map[string]*Chain{}}
	fw.Tables[name] = t
	return t
}

func ensureChain(t *Table, name string, ctype ChainType, policy Verdict) *Chain {
	if ch, ok := t.Chains[name]; ok {
		// If chain already exists as user, but now declared builtin/policy, update.
		if ctype == ChainBuiltin {
			ch.Type = ChainBuiltin
			if policy != "" {
				ch.Policy = policy
			}
		}
		return ch
	}
	ch := &Chain{Name: name, Type: ctype, Policy: policy}
	t.Chains[name] = ch
	return ch
}

func ensureBuiltinDefaults(fw *Firewall) {
	t := ensureTable(fw, TableFilter)
	// If not declared, assume typical defaults (policy ACCEPT commonly, but iptables-restore style expects explicit).
	// We'll create them as builtin with policy ACCEPT unless already set.
	defs := []string{"INPUT", "FORWARD", "OUTPUT"}
	for _, c := range defs {
		if _, ok := t.Chains[c]; !ok {
			ensureChain(t, c, ChainBuiltin, VAccept)
		} else if t.Chains[c].Type != ChainBuiltin {
			// If user created chain with same name, still treat as builtin for simulation safety.
			t.Chains[c].Type = ChainBuiltin
			if t.Chains[c].Policy == "" {
				t.Chains[c].Policy = VAccept
			}
		}
	}
}

func parseChainDef(t *Table, line string) error {
	// :CHAIN POLICY [0:0]
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return fmt.Errorf("bad chain def: %q", line)
	}
	name := strings.TrimPrefix(parts[0], ":")
	pol := strings.ToUpper(parts[1])

	ctype := ChainUser
	var policy Verdict
	if pol != "-" {
		ctype = ChainBuiltin
		switch pol {
		case "ACCEPT":
			policy = VAccept
		case "DROP":
			policy = VDrop
		case "REJECT":
			policy = VReject
		default:
			return fmt.Errorf("unsupported policy: %s", pol)
		}
	}
	ensureChain(t, name, ctype, policy)
	return nil
}

func parsePolicyLine(t *Table, line string) error {
	// -P INPUT DROP
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return fmt.Errorf("bad -P line: %q", line)
	}
	ch := parts[1]
	p := strings.ToUpper(parts[2])
	var pol Verdict
	switch p {
	case "ACCEPT":
		pol = VAccept
	case "DROP":
		pol = VDrop
	case "REJECT":
		pol = VReject
	default:
		return fmt.Errorf("unsupported policy: %s", p)
	}
	ensureChain(t, ch, ChainBuiltin, pol)
	return nil
}

func parseRule(lineNo int, line string) (*Rule, string, error) {
	toks, err := shellSplit(line)
	if err != nil {
		return nil, "", err
	}
	if len(toks) < 3 {
		return nil, "", fmt.Errorf("short rule: %q", line)
	}
	// -A CHAIN ...
	if toks[0] != "-A" && toks[0] != "-I" {
		return nil, "", fmt.Errorf("expected -A/-I, got %s", toks[0])
	}
	chainName := toks[1]
	args := toks[2:]

	r := &Rule{LineNo: lineNo, Raw: line}
	r.Match = Match{
		CTStates: map[string]bool{},
	}
	r.Target = Target{Kind: TUnknown, Raw: ""}

	i := 0
	for i < len(args) {
		a := args[i]

		switch a {
		case "-p", "--protocol":
			i++
			if i >= len(args) {
				return nil, "", fmt.Errorf("missing proto after %s", a)
			}
			r.Match.Proto = strings.ToLower(args[i])

		case "-s", "--source":
			i++
			if i >= len(args) {
				return nil, "", fmt.Errorf("missing source after %s", a)
			}
			_, n, err := net.ParseCIDR(args[i])
			if err != nil {
				// allow plain IP
				ip := net.ParseIP(args[i])
				if ip == nil {
					return nil, "", fmt.Errorf("bad source: %s", args[i])
				}
				n = &net.IPNet{IP: ip, Mask: net.CIDRMask(32, 32)}
			}
			r.Match.SrcCIDR = n

		case "-d", "--destination":
			i++
			if i >= len(args) {
				return nil, "", fmt.Errorf("missing destination after %s", a)
			}
			_, n, err := net.ParseCIDR(args[i])
			if err != nil {
				ip := net.ParseIP(args[i])
				if ip == nil {
					return nil, "", fmt.Errorf("bad destination: %s", args[i])
				}
				n = &net.IPNet{IP: ip, Mask: net.CIDRMask(32, 32)}
			}
			r.Match.DstCIDR = n

		case "-i", "--in-interface":
			i++
			if i >= len(args) {
				return nil, "", fmt.Errorf("missing iface after %s", a)
			}
			r.Match.InIface = args[i]

		case "-o", "--out-interface":
			i++
			if i >= len(args) {
				return nil, "", fmt.Errorf("missing iface after %s", a)
			}
			r.Match.OutIface = args[i]

		case "--sport", "--source-port":
			i++
			if i >= len(args) {
				return nil, "", fmt.Errorf("missing sport after %s", a)
			}
			ps, err := parsePortSpec(args[i])
			if err != nil {
				return nil, "", err
			}
			r.Match.SrcPort = ps

		case "--dport", "--destination-port":
			i++
			if i >= len(args) {
				return nil, "", fmt.Errorf("missing dport after %s", a)
			}
			ps, err := parsePortSpec(args[i])
			if err != nil {
				return nil, "", err
			}
			r.Match.DstPort = ps

		case "-m":
			// module, we parse some next options accordingly; just skip module token itself
			i++
			if i >= len(args) {
				return nil, "", fmt.Errorf("missing module name after -m")
			}
			mod := strings.ToLower(args[i])
			_ = mod // handled by subsequent options

		case "--ctstate":
			i++
			if i >= len(args) {
				return nil, "", fmt.Errorf("missing states after --ctstate")
			}
			for _, st := range strings.Split(args[i], ",") {
				s := strings.ToUpper(strings.TrimSpace(st))
				if s != "" {
					r.Match.CTStates[s] = true
				}
			}

		case "--mark":
			// iptables mark match syntax can be value[/mask]
			i++
			if i >= len(args) {
				return nil, "", fmt.Errorf("missing value after --mark")
			}
			val, mask, err := parseMarkMatch(args[i])
			if err != nil {
				return nil, "", err
			}
			r.Match.HasMarkMatch = true
			r.Match.MarkValue = val
			r.Match.MarkMask = mask

		case "--tcp-flags":
			// --tcp-flags mask comp
			i++
			if i >= len(args) {
				return nil, "", fmt.Errorf("missing mask after --tcp-flags")
			}
			maskSet := parseFlagSet(args[i])
			i++
			if i >= len(args) {
				return nil, "", fmt.Errorf("missing comp after --tcp-flags <mask>")
			}
			compSet := parseFlagSet(args[i])
			r.Match.TCPFlags = &TCPFlagsSpec{Mask: maskSet, Comp: compSet}

		case "--comment":
			i++
			if i >= len(args) {
				return nil, "", fmt.Errorf("missing comment after --comment")
			}
			r.Match.Comment = args[i]

		case "-j", "--jump":
			i++
			if i >= len(args) {
				return nil, "", fmt.Errorf("missing target after %s", a)
			}
			tgt := strings.ToUpper(args[i])
			t, err := parseTargetJump(tgt, args, &i)
			if err != nil {
				return nil, "", err
			}
			r.Target = t

		case "-g", "--goto":
			i++
			if i >= len(args) {
				return nil, "", fmt.Errorf("missing chain after %s", a)
			}
			ch := args[i]
			r.Target = Target{Kind: TGoto, ChainName: ch, Raw: "GOTO " + ch}

		case "--log-prefix":
			// used after -j LOG; accept even if order differs
			i++
			if i >= len(args) {
				return nil, "", fmt.Errorf("missing value after --log-prefix")
			}
			// attach to rule target if LOG, else store in raw
			if r.Target.Kind == TLog {
				r.Target.LogPrefix = args[i]
			}

		case "--set-mark":
			// MARK target option
			i++
			if i >= len(args) {
				return nil, "", fmt.Errorf("missing value after --set-mark")
			}
			val, mask, err := parseSetMark(args[i])
			if err != nil {
				return nil, "", err
			}
			// If target not MARK yet, set it (support "iptables -j MARK --set-mark ...")
			if r.Target.Kind == TUnknown {
				r.Target = Target{Kind: TMark, SetMark: val, SetMarkMask: mask, Raw: "MARK"}
			} else if r.Target.Kind == TMark {
				r.Target.SetMark = val
				r.Target.SetMarkMask = mask
			} else {
				// ignore in non-MARK target
			}

		default:
			// unknown tokens: ignore for now (you can extend)
		}

		i++
	}

	if r.Target.Kind == TUnknown {
		// rule with no action: increments counters only
		r.Target = Target{Kind: TUnknown, Raw: "NONE"}
	}

	return r, chainName, nil
}

func parseTargetJump(tgt string, args []string, i *int) (Target, error) {
	switch tgt {
	case "ACCEPT":
		return Target{Kind: TAccept, Raw: "ACCEPT"}, nil
	case "DROP":
		return Target{Kind: TDrop, Raw: "DROP"}, nil
	case "REJECT":
		return Target{Kind: TReject, Raw: "REJECT"}, nil
	case "RETURN":
		return Target{Kind: TReturn, Raw: "RETURN"}, nil
	case "LOG":
		// may have --log-prefix later
		return Target{Kind: TLog, Raw: "LOG"}, nil
	case "MARK":
		// MARK expects --set-mark <val[/mask]> later
		return Target{Kind: TMark, Raw: "MARK"}, nil
	default:
		// user-defined chain jump
		return Target{Kind: TJump, ChainName: tgt, Raw: "JUMP " + tgt}, nil
	}
}

func parsePortSpec(s string) (*PortSpec, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "0" {
		return &PortSpec{Any: true}, nil
	}
	if strings.Contains(s, ":") {
		parts := strings.SplitN(s, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("bad port range: %s", s)
		}
		from := 0
		to := 65535
		if parts[0] != "" {
			v, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("bad port: %s", parts[0])
			}
			from = v
		}
		if parts[1] != "" {
			v, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("bad port: %s", parts[1])
			}
			to = v
		}
		return &PortSpec{Range: true, From: from, To: to}, nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return nil, fmt.Errorf("bad port: %s", s)
	}
	return &PortSpec{Single: v}, nil
}

func parseMarkMatch(s string) (val uint32, mask uint32, err error) {
	// value[/mask]
	parts := strings.SplitN(s, "/", 2)
	v, err := parseUint32(parts[0])
	if err != nil {
		return 0, 0, err
	}
	m := uint32(0xFFFFFFFF)
	if len(parts) == 2 {
		mm, err := parseUint32(parts[1])
		if err != nil {
			return 0, 0, err
		}
		m = mm
	}
	return v, m, nil
}

func parseSetMark(s string) (val uint32, mask uint32, err error) {
	// iptables MARK --set-mark value[/mask]
	parts := strings.SplitN(s, "/", 2)
	v, err := parseUint32(parts[0])
	if err != nil {
		return 0, 0, err
	}
	m := uint32(0xFFFFFFFF)
	if len(parts) == 2 {
		mm, err := parseUint32(parts[1])
		if err != nil {
			return 0, 0, err
		}
		m = mm
	}
	return v, m, nil
}

func parseUint32(s string) (uint32, error) {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		u, err := strconv.ParseUint(s[2:], 16, 32)
		return uint32(u), err
	}
	u, err := strconv.ParseUint(s, 10, 32)
	return uint32(u), err
}

func parseIPLoose(s string) (net.IP, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, errors.New("empty")
	}
	// allow "1.2.3.4/24" but take IP portion
	if strings.Contains(s, "/") {
		ipStr := strings.SplitN(s, "/", 2)[0]
		ip := net.ParseIP(ipStr)
		if ip == nil {
			return nil, fmt.Errorf("bad ip: %s", s)
		}
		return ip, nil
	}
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("bad ip: %s", s)
	}
	return ip, nil
}

func parseFlagSet(s string) map[string]bool {
	out := map[string]bool{}
	for _, p := range strings.Split(s, ",") {
		f := strings.ToUpper(strings.TrimSpace(p))
		if f != "" {
			out[f] = true
		}
	}
	return out
}

// Very small shell-like splitter supporting quotes.
func shellSplit(s string) ([]string, error) {
	var out []string
	var cur strings.Builder
	inQuote := byte(0)
	esc := false

	flush := func() {
		if cur.Len() > 0 {
			out = append(out, cur.String())
			cur.Reset()
		}
	}

	for i := 0; i < len(s); i++ {
		c := s[i]

		if esc {
			cur.WriteByte(c)
			esc = false
			continue
		}
		if c == '\\' {
			esc = true
			continue
		}

		if inQuote != 0 {
			if c == inQuote {
				inQuote = 0
			} else {
				cur.WriteByte(c)
			}
			continue
		}

		if c == '"' || c == '\'' {
			inQuote = c
			continue
		}

		if c == ' ' || c == '\t' {
			flush()
			continue
		}
		cur.WriteByte(c)
	}
	if esc {
		return nil, fmt.Errorf("dangling escape")
	}
	if inQuote != 0 {
		return nil, fmt.Errorf("unterminated quote")
	}
	flush()
	return out, nil
}

// ------------------------------ Simulation Engine ------------------------------

type frame struct {
	table string
	chain string
	rule  int // next rule index upon return
}

func (fw *Firewall) Simulate(table TableName, startChain string, pkt Packet) (*SimResult, error) {
	t, ok := fw.Tables[table]
	if !ok {
		return nil, fmt.Errorf("table not found: %s", table)
	}
	ch, ok := t.Chains[startChain]
	if !ok {
		return nil, fmt.Errorf("chain not found: %s", startChain)
	}

	res := &SimResult{
		Table:     string(table),
		Chain:     startChain,
		Verdict:   VNone,
		Trace:     []string{},
		Logs:      []string{},
		FinalMark: pkt.Mark,
	}

	// Call stack for -j semantics.
	var stack []frame

	curChain := ch
	curChainName := startChain
	curRule := 0

	// This "caller chain" concept is needed for -g RETURN semantics.
	// With -g, we do not push a return address. RETURN then returns to
	// the chain that "called us via -j" (or top-level end).
	// We can model this by: if RETURN happens and stack is empty, we exit
	// the current traversal (end of top-level chain). If stack not empty,
	// pop and continue from caller frame.
	for {
		if curRule >= len(curChain.Rules) {
			// End of chain:
			if curChain.Type == ChainBuiltin {
				// builtin chain: apply policy
				res.Verdict = curChain.Policy
				res.Reason = fmt.Sprintf("end of builtin chain %s: policy %s", curChainName, curChain.Policy)
				res.FinalMark = pkt.Mark
				return res, nil
			}
			// user-defined chain: implicit RETURN
			res.Trace = append(res.Trace, fmt.Sprintf("end of user chain %s => implicit RETURN", curChainName))
			ok := doReturn(&stack, &curChainName, &curChain, &curRule, t)
			if !ok {
				// no caller: returning from top-level chain means "no verdict" => builtin policy should handle,
				// but top-level should always be builtin in practical usage.
				// We'll treat as ACCEPT by default if start chain is user (rare).
				res.Verdict = VAccept
				res.Reason = "returned from top-level user chain: default ACCEPT"
				res.FinalMark = pkt.Mark
				return res, nil
			}
			continue
		}

		r := curChain.Rules[curRule]
		// Even if rule doesn't match, increment index only.
		if !ruleMatches(r, &pkt) {
			curRule++
			continue
		}

		// matched => increment counters
		r.Counters[0]++
		r.Counters[1] += pkt.Bytes
		curChain.Counters[0]++
		curChain.Counters[1] += pkt.Bytes

		res.Trace = append(res.Trace, fmt.Sprintf("%s[%d] MATCH line %d: %s", curChainName, curRule, r.LineNo, r.Raw))

		verdict, err := applyTarget(r.Target, &pkt, &stack, t, &curChainName, &curChain, &curRule, res)
		if err != nil {
			return nil, err
		}

		switch verdict {
		case VNone:
			// continue
			curRule++
		case VReturn:
			res.Trace = append(res.Trace, fmt.Sprintf("%s[%d] => RETURN", curChainName, curRule))
			ok := doReturn(&stack, &curChainName, &curChain, &curRule, t)
			if !ok {
				// returning out of top-level chain => apply builtin policy if top-level is builtin, else default ACCEPT.
				if start, ok2 := t.Chains[startChain]; ok2 && start.Type == ChainBuiltin {
					res.Verdict = start.Policy
					res.Reason = fmt.Sprintf("RETURN from top-level builtin chain %s => policy %s", startChain, start.Policy)
				} else {
					res.Verdict = VAccept
					res.Reason = "RETURN from top-level user chain => default ACCEPT"
				}
				res.FinalMark = pkt.Mark
				return res, nil
			}
		case VAccept, VDrop, VReject:
			res.Verdict = verdict
			res.Reason = fmt.Sprintf("decided by chain %s rule %d (line %d)", curChainName, curRule, r.LineNo)
			res.FinalMark = pkt.Mark
			return res, nil
		default:
			// should not happen
			return nil, fmt.Errorf("unknown verdict: %v", verdict)
		}
	}
}

func doReturn(stack *[]frame, curChainName *string, curChain **Chain, curRule *int, t *Table) (ok bool) {
	if len(*stack) == 0 {
		return false
	}
	// Pop and restore caller.
	top := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]

	ch, ok := t.Chains[top.chain]
	if !ok {
		// should not happen
		return false
	}
	*curChainName = top.chain
	*curChain = ch
	*curRule = top.rule
	return true
}

func applyTarget(tgt Target, pkt *Packet, stack *[]frame, table *Table, curChainName *string, curChain **Chain, curRule *int, res *SimResult) (Verdict, error) {
	switch tgt.Kind {
	case TUnknown:
		// no effect
		return VNone, nil
	case TAccept:
		return VAccept, nil
	case TDrop:
		return VDrop, nil
	case TReject:
		return VReject, nil
	case TReturn:
		return VReturn, nil
	case TLog:
		pfx := tgt.LogPrefix
		if pfx == "" {
			pfx = "LOG "
		}
		msg := fmt.Sprintf("%sproto=%s src=%s dst=%s sport=%d dport=%d in=%s out=%s mark=0x%X ct=%s",
			pfx,
			pkt.Proto,
			pkt.SrcIP.String(),
			pkt.DstIP.String(),
			pkt.SrcPort,
			pkt.DstPort,
			pkt.InIface,
			pkt.OutIface,
			pkt.Mark,
			pkt.CTState,
		)
		res.Logs = append(res.Logs, msg)
		return VNone, nil
	case TMark:
		// Simulate MARK --set-mark value[/mask]
		mask := tgt.SetMarkMask
		if mask == 0 {
			mask = 0xFFFFFFFF
		}
		// iptables MARK semantics: (mark & ~mask) | (value & mask)
		old := pkt.Mark
		pkt.Mark = (pkt.Mark & ^mask) | (tgt.SetMark & mask)
		res.Trace = append(res.Trace, fmt.Sprintf("MARK: 0x%X -> 0x%X (value=0x%X mask=0x%X)", old, pkt.Mark, tgt.SetMark, mask))
		return VNone, nil
	case TJump:
		// push return frame and enter chain
		next := tgt.ChainName
		ch, ok := table.Chains[next]
		if !ok {
			return VDrop, fmt.Errorf("jump to unknown chain: %s", next)
		}
		*stack = append(*stack, frame{table: string(table.Name), chain: *curChainName, rule: *curRule + 1})
		res.Trace = append(res.Trace, fmt.Sprintf("JUMP: %s -> %s (return to %s[%d])", *curChainName, next, *curChainName, *curRule+1))
		*curChainName = next
		*curChain = ch
		*curRule = 0
		return VNone, nil
	case TGoto:
		// no push, just transfer control
		next := tgt.ChainName
		ch, ok := table.Chains[next]
		if !ok {
			return VDrop, fmt.Errorf("goto unknown chain: %s", next)
		}
		res.Trace = append(res.Trace, fmt.Sprintf("GOTO: %s -> %s (no return point)", *curChainName, next))
		*curChainName = next
		*curChain = ch
		*curRule = 0
		return VNone, nil
	default:
		return VDrop, fmt.Errorf("unsupported target: %v", tgt.Kind)
	}
}

func ruleMatches(r *Rule, pkt *Packet) bool {
	m := r.Match

	// proto
	if m.Proto != "" && m.Proto != "any" {
		if strings.ToLower(pkt.Proto) != strings.ToLower(m.Proto) {
			return false
		}
	}

	// src/dst
	if m.SrcCIDR != nil {
		if pkt.SrcIP == nil || !m.SrcCIDR.Contains(pkt.SrcIP) {
			return false
		}
	}
	if m.DstCIDR != nil {
		if pkt.DstIP == nil || !m.DstCIDR.Contains(pkt.DstIP) {
			return false
		}
	}

	// iface
	if m.InIface != "" && m.InIface != pkt.InIface {
		return false
	}
	if m.OutIface != "" && m.OutIface != pkt.OutIface {
		return false
	}

	// ports (tcp/udp typically)
	if m.SrcPort != nil && !portMatch(m.SrcPort, pkt.SrcPort) {
		return false
	}
	if m.DstPort != nil && !portMatch(m.DstPort, pkt.DstPort) {
		return false
	}

	// ctstate
	if len(m.CTStates) > 0 {
		if pkt.CTState == "" {
			return false
		}
		if !m.CTStates[strings.ToUpper(pkt.CTState)] {
			return false
		}
	}

	// mark
	if m.HasMarkMatch {
		// match if (mark & mask) == (value & mask)
		if (pkt.Mark & m.MarkMask) != (m.MarkValue & m.MarkMask) {
			return false
		}
	}

	// tcp flags
	if m.TCPFlags != nil {
		// Only check flags in mask set; they must equal comp set.
		// Meaning: for each flag in mask, pktFlag must be 1 iff comp has it.
		for f := range m.TCPFlags.Mask {
			present := pkt.TCPFlags[strings.ToUpper(f)]
			want := m.TCPFlags.Comp[strings.ToUpper(f)]
			if present != want {
				return false
			}
		}
	}

	return true
}

func portMatch(ps *PortSpec, p int) bool {
	if ps == nil {
		return true
	}
	if ps.Any {
		return true
	}
	if ps.Range {
		return p >= ps.From && p <= ps.To
	}
	return p == ps.Single
}

// ------------------------------ Printing ------------------------------

func printFirewall(fw *Firewall) {
	var tables []string
	for tn := range fw.Tables {
		tables = append(tables, string(tn))
	}
	sort.Strings(tables)

	for _, tname := range tables {
		t := fw.Tables[TableName(tname)]
		fmt.Printf("*%s\n", t.Name)

		var chains []string
		for cn := range t.Chains {
			chains = append(chains, cn)
		}
		sort.Strings(chains)

		for _, cn := range chains {
			ch := t.Chains[cn]
			if ch.Type == ChainBuiltin {
				fmt.Printf(":%s %s [%d:%d]\n", ch.Name, ch.Policy, ch.Counters[0], ch.Counters[1])
			} else {
				fmt.Printf(":%s - [%d:%d]\n", ch.Name, ch.Counters[0], ch.Counters[1])
			}
		}

		for _, cn := range chains {
			ch := t.Chains[cn]
			if len(ch.Rules) == 0 {
				continue
			}
			fmt.Printf("\n# Chain %s (%s) rules=%d\n", ch.Name, chainTypeStr(ch.Type), len(ch.Rules))
			for i, r := range ch.Rules {
				fmt.Printf("[%d] line %d  %s\n", i, r.LineNo, r.Raw)
			}
		}

		fmt.Printf("COMMIT\n\n")
	}
}

func chainTypeStr(ct ChainType) string {
	if ct == ChainBuiltin {
		return "builtin"
	}
	return "user"
}

// ------------------------------ Extra Helpers ------------------------------

var (
	reCommaSpace = regexp.MustCompile(`\s*,\s*`)
)
