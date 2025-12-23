package iptable

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/k0kubun/pp"
)

type tableType int

const (
	raw      tableType = iota // Before conntrack
	mangle   tableType = iota // To modify packet (TTL, marks, Qos)
	nat      tableType = iota // To SNAT/DNAT
	filter   tableType = iota // To filter out
	security tableType = iota // SELinux
)

type chainType int

const (
	prerouting  chainType = iota
	input       chainType = iota
	forward     chainType = iota
	output      chainType = iota
	postrouting chainType = iota
)

var rawBuiltinChains = [...]chainType{prerouting, output}
var mangleBuiltinChains = [...]chainType{prerouting, input, forward, output, postrouting}
var natBuiltinChains = [...]chainType{prerouting, input, forward, output}
var filterBuiltinChains = [...]chainType{input, forward, output}
var securityBuiltinChains = [...]chainType{input, forward, output}

var inboundChaining = [...]chainType{prerouting, input}
var outboundChaining = [...]chainType{output, postrouting}
var forwardChaining = [...]chainType{prerouting, forward, postrouting}

type iptableBackend struct{}

func NewBackend() *iptableBackend {
	return &iptableBackend{}
}

type table struct {
	Name   string
	Chains map[string]*chain
}

type chain struct {
	Name   string
	Rules  []rule
	Policy string
}

type rule struct {
	Raw        string
	Chain      string
	JumpTarget string
}

type counter struct {
	Packets uint64
	Bytes   uint64
}

func (b *iptableBackend) ListChains(_ string) ([]string, error) {
	config, err := b.fetch()
	if err != nil {
		return nil, err
	}
	pp.Println(config)
	return []string{}, nil
}

func (b *iptableBackend) fetch() (map[string]table, error) {
	output, err := b.runProces()
	if err != nil {
		return nil, err
	}
	return b.parseTables(output)
}

func (b *iptableBackend) parseTables(input string) (map[string]table, error) {
	tables := make(map[string]table)

	lines := strings.Split(input, "\n")
	var currentTable table

	for _, line := range lines {
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "*"):
			name := strings.TrimPrefix(line, "*")
			currentTable = table{
				Name:   name,
				Chains: map[string]*chain{},
			}
			tables[name] = currentTable
		case line == "COMMIT":
			currentTable = table{}
		case strings.HasPrefix(line, ":"):
			fmt.Println(line)
			chainItem, _ := b.parseChain(line)
			currentTable.Chains[chainItem.Name] = &chainItem
		case strings.HasPrefix(line, "-A"):
			ruleItem, err := b.parseRule(line)
			if err != nil {
				return nil, err
			}
			currentTable.Chains[ruleItem.Chain].Rules = append(currentTable.Chains[ruleItem.Chain].Rules, ruleItem)
		}
	}

	return tables, nil
}

func (b *iptableBackend) parseRule(input string) (rule, error) {
	parts := strings.Fields(input)
	if len(parts) < 3 {
		return rule{}, errors.New("invalid iptables rule")
	}
	chainName := parts[1]
	ruleItem := rule{
		Chain: chainName,
		Raw:   input,
	}

	for i := 2; i < len(parts); i++ {
		switch parts[i] {
		case "-j":
			if i+1 < len(parts) {
				ruleItem.JumpTarget = parts[i+1]
				i++
			}
		}
	}

	return ruleItem, nil
}

func (b *iptableBackend) parseChain(line string) (chain, error) {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return chain{}, errors.New("invalid iptable chain format")
	}

	name := strings.TrimPrefix(parts[0], ":")

	chain := chain{
		Name:   name,
		Rules:  []rule{},
		Policy: parts[1],
	}

	return chain, nil
}

func (b *iptableBackend) parseCounter(raw string) counter {
	// [packets:bytes]
	raw = strings.Trim(raw, "[]")
	parts := strings.Split(raw, ":")
	if len(parts) != 2 {
		return counter{}
	}
	fmt.Println(parts)
	packets, _ := strconv.ParseUint(parts[0], 10, 64)
	bytes, _ := strconv.ParseUint(parts[1], 10, 64)

	return counter{Packets: packets, Bytes: bytes}
}

func (b *iptableBackend) runProces() (string, error) {
	// -c add the counters , "-c"
	cmd := exec.Command("iptables-save")
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err

	if errRun := cmd.Run(); errRun != nil {
		fmt.Println(errRun.Error())
		return "", errors.New(err.String())
	}

	fmt.Println(out.String())

	return out.String(), nil
}
