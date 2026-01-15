//go:build IPTABLE_BINARY

package iptable

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

type iptableBackend struct {
	stdout string
}

func NewBackend() *iptableBackend {
	return &iptableBackend{}
}

func (b *iptableBackend) ListChains(_ string) (map[string]Table, error) {
	config, err := b.Fetch()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (b *iptableBackend) Fetch() (map[string]Table, error) {
	output, err := b.runProces()
	b.stdout = output
	if err != nil {
		return nil, err
	}
	return b.parseTables(output)
}

func (b *iptableBackend) GetStdout() string {
	return b.stdout
}

func (b *iptableBackend) parseTables(input string) (map[string]Table, error) {
	tables := make(map[string]Table)

	lines := strings.Split(input, "\n")
	var currentTable Table

	for _, line := range lines {
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "*"):
			name := strings.TrimPrefix(line, "*")
			currentTable = Table{
				Name:   name,
				Chains: map[string]*chain{},
			}
			tables[name] = currentTable
		case line == "COMMIT":
			currentTable = Table{}
		case strings.HasPrefix(line, ":"):
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

func (b *iptableBackend) runProces() (string, error) {
	// -c add the counters , "-c"
	cmd := exec.Command("iptables-save")
	var out bytes.Buffer
	var err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err

	if errRun := cmd.Run(); errRun != nil {
		return "", errors.New(err.String())
	}

	return out.String(), nil
}
