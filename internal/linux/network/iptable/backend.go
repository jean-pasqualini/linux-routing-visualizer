// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package iptable

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

// Adds the output of stderr to exec.ExitError
type Error struct {
	exec.ExitError
	cmd        exec.Cmd
	msg        string
	exitStatus *int //for overriding
}

func (e *Error) ExitStatus() int {
	if e.exitStatus != nil {
		return *e.exitStatus
	}
	return e.Sys().(syscall.WaitStatus).ExitStatus()
}

func (e *Error) Error() string {
	return fmt.Sprintf("running %v: exit status %v: %v", e.cmd.Args, e.ExitStatus(), e.msg)
}

var isNotExistPatterns = []string{
	"Bad rule (does a matching rule exist in that chain?).\n",
	"No chain/target/match by that name.\n",
	"No such file or directory",
	"does not exist",
}

// IsNotExist returns true if the error is due to the chain or rule not existing
func (e *Error) IsNotExist() bool {
	for _, str := range isNotExistPatterns {
		if strings.Contains(e.msg, str) {
			return true
		}
	}
	return false
}

// Protocol to differentiate between IPv4 and IPv6
type Protocol byte

const (
	ProtocolIPv4 Protocol = iota
	ProtocolIPv6
)

type IPTables struct {
	path              string
	proto             Protocol
	hasCheck          bool
	hasWait           bool
	waitSupportSecond bool
	hasRandomFully    bool
	v1                int
	v2                int
	v3                int
	mode              string // the underlying iptables operating mode, e.g. nf_tables
	timeout           int    // time to wait for the iptables lock, default waits forever
}

// Stat represents a structured statistic entry.
type Stat struct {
	Packets     uint64     `json:"pkts"`
	Bytes       uint64     `json:"bytes"`
	Target      string     `json:"target"`
	Protocol    string     `json:"prot"`
	Opt         string     `json:"opt"`
	Input       string     `json:"in"`
	Output      string     `json:"out"`
	Source      *net.IPNet `json:"source"`
	Destination *net.IPNet `json:"destination"`
	Options     string     `json:"options"`
}

type option func(*IPTables)

func NewBackend(opts ...option) (*IPTables, error) {

	ipt := &IPTables{
		proto:   ProtocolIPv4,
		timeout: 0,
		path:    "",
	}

	for _, opt := range opts {
		opt(ipt)
	}

	// if path wasn't preset through New(Path()), autodiscover it
	cmd := ""
	if ipt.path == "" {
		cmd = getIptablesCommand(ipt.proto)
	} else {
		cmd = ipt.path
	}
	path, err := exec.LookPath(cmd)
	if err != nil {
		return nil, err
	}
	ipt.path = path

	vstring, err := getIptablesVersionString(path)
	if err != nil {
		return nil, fmt.Errorf("could not get iptables version: %v", err)
	}
	v1, v2, v3, mode, err := extractIptablesVersion(vstring)
	if err != nil {
		return nil, fmt.Errorf("failed to extract iptables version from [%s]: %v", vstring, err)
	}
	ipt.v1 = v1
	ipt.v2 = v2
	ipt.v3 = v3
	ipt.mode = mode

	return ipt, nil
}

// Proto returns the protocol used by this IPTables.
func (ipt *IPTables) Proto() Protocol {
	return ipt.proto
}

// List rules in specified table/chain
func (ipt *IPTables) ListById(table, chain string, id int) (string, error) {
	args := []string{"-t", table, "-S", chain, strconv.Itoa(id)}
	rule, err := ipt.executeList(args)
	if err != nil {
		return "", err
	}
	return rule[0], nil
}

// List rules in specified table/chain
func (ipt *IPTables) List(table, chain string) ([]string, error) {
	args := []string{"-t", table, "-S", chain}
	return ipt.executeList(args)
}

// List rules (with counters) in specified table/chain
func (ipt *IPTables) ListWithCounters(table, chain string) ([]string, error) {
	args := []string{"-t", table, "-v", "-S", chain}
	return ipt.executeList(args)
}

// ListChains returns a slice containing the name of each chain in the specified table.
func (ipt *IPTables) ListChains(table string) ([]string, error) {
	args := []string{"-t", table, "-S"}

	result, err := ipt.executeList(args)
	if err != nil {
		return nil, err
	}

	// Iterate over rules to find all default (-P) and user-specified (-N) chains.
	// Chains definition always come before rules.
	// Format is the following:
	// -P OUTPUT ACCEPT
	// -N Custom
	var chains []string
	for _, val := range result {
		if strings.HasPrefix(val, "-P") || strings.HasPrefix(val, "-N") {
			chains = append(chains, strings.Fields(val)[1])
		} else {
			break
		}
	}
	return chains, nil
}

func (ipt *IPTables) executeList(args []string) ([]string, error) {
	var stdout bytes.Buffer
	if err := ipt.runWithOutput(args, &stdout); err != nil {
		return nil, err
	}

	rules := strings.Split(stdout.String(), "\n")

	// strip trailing newline
	if len(rules) > 0 && rules[len(rules)-1] == "" {
		rules = rules[:len(rules)-1]
	}

	for i, rule := range rules {
		rules[i] = filterRuleOutput(rule)
	}

	return rules, nil
}

// Return version components of the underlying iptables command
func (ipt *IPTables) GetIptablesVersion() (int, int, int) {
	return ipt.v1, ipt.v2, ipt.v3
}

// run runs an iptables command with the given arguments, ignoring
// any stdout output
func (ipt *IPTables) run(args ...string) error {
	return ipt.runWithOutput(args, nil)
}

// runWithOutput runs an iptables command with the given arguments,
// writing any stdout output to the given writer
func (ipt *IPTables) runWithOutput(args []string, stdout io.Writer) error {
	args = append([]string{"sudo", ipt.path}, args...)

	var stderr bytes.Buffer
	cmd := exec.Cmd{
		Path:   ipt.path,
		Args:   args,
		Stdout: stdout,
		Stderr: &stderr,
	}

	fmt.Println(cmd)
	if err := cmd.Run(); err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			return &Error{*e, cmd, stderr.String(), nil}
		default:
			return err
		}
	}

	return nil
}

// getIptablesCommand returns the correct command for the given protocol, either "iptables" or "ip6tables".
func getIptablesCommand(proto Protocol) string {
	if proto == ProtocolIPv6 {
		return "ip6tables"
	} else {
		return "iptables"
	}
}

// getIptablesVersion returns the first three components of the iptables version
// and the operating mode (e.g. nf_tables or legacy)
// e.g. "iptables v1.3.66" would return (1, 3, 66, legacy, nil)
func extractIptablesVersion(str string) (int, int, int, string, error) {
	versionMatcher := regexp.MustCompile(`v([0-9]+)\.([0-9]+)\.([0-9]+)(?:\s+\((\w+))?`)
	result := versionMatcher.FindStringSubmatch(str)
	if result == nil {
		return 0, 0, 0, "", fmt.Errorf("no iptables version found in string: %s", str)
	}

	v1, err := strconv.Atoi(result[1])
	if err != nil {
		return 0, 0, 0, "", err
	}

	v2, err := strconv.Atoi(result[2])
	if err != nil {
		return 0, 0, 0, "", err
	}

	v3, err := strconv.Atoi(result[3])
	if err != nil {
		return 0, 0, 0, "", err
	}

	mode := "legacy"
	if result[4] != "" {
		mode = result[4]
	}
	return v1, v2, v3, mode, nil
}

// Runs "iptables --version" to get the version string
func getIptablesVersionString(path string) (string, error) {
	cmd := exec.Command(path, "--version")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// Checks if a rule specification exists for a table
func (ipt *IPTables) existsForOldIptables(table, chain string, rulespec []string) (bool, error) {
	rs := strings.Join(append([]string{"-A", chain}, rulespec...), " ")
	args := []string{"-t", table, "-S"}
	var stdout bytes.Buffer
	err := ipt.runWithOutput(args, &stdout)
	if err != nil {
		return false, err
	}
	return strings.Contains(stdout.String(), rs), nil
}

// counterRegex is the regex used to detect nftables counter format
var counterRegex = regexp.MustCompile(`^\[([0-9]+):([0-9]+)\] `)

// filterRuleOutput works around some inconsistencies in output.
// For example, when iptables is in legacy vs. nftables mode, it produces
// different results.
func filterRuleOutput(rule string) string {
	out := rule

	// work around an output difference in nftables mode where counters
	// are output in iptables-save format, rather than iptables -S format
	// The string begins with "[0:0]"
	//
	// Fixes #49
	if groups := counterRegex.FindStringSubmatch(out); groups != nil {
		// drop the brackets
		out = out[len(groups[0]):]
		out = fmt.Sprintf("%s -c %s %s", out, groups[1], groups[2])
	}

	return out
}
