package iptable

import (
	"errors"
	"strconv"
	"strings"
)

func (b *iptableBackend) parseChain(line string) (chain, error) {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return chain{}, errors.New("invalid iptable chain format")
	}

	name := strings.TrimPrefix(parts[0], ":")

	chain := chain{
		Name:   name,
		Rules:  []Rule{},
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
	packets, _ := strconv.ParseUint(parts[0], 10, 64)
	bytes, _ := strconv.ParseUint(parts[1], 10, 64)

	return counter{Packets: packets, Bytes: bytes}
}
