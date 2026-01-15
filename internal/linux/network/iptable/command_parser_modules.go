package iptable

import "strings"

func (b *iptableBackend) parseRuleModuleTCPUDP(input string, ruleItem Rule) Rule {
	parts := strings.Fields(input)

	for i := 2; i < len(parts); i++ {
		switch parts[i] {
		case "--dport":
			if i+1 < len(parts) {
				ruleItem.Filter.To.Port = parts[i+1]
				i++
			}
		case "--sport":
			if i+1 < len(parts) {
				ruleItem.Filter.From.Port = parts[i+1]
				i++
			}
		}
	}

	return ruleItem
}

func (b *iptableBackend) parseRuleModuleConntrack(input string, ruleItem Rule) Rule {
	parts := strings.Fields(input)

	for i := 2; i < len(parts); i++ {
		switch parts[i] {
		case "--ctstate":
			if i+1 < len(parts) {
				states := strings.Split(parts[i+1], ",")
				ruleItem.Filter.ConnectionState = states
				i++
			}
		}
	}

	return ruleItem
}

func (b *iptableBackend) parseRuleModuleAddrType(input string, ruleItem Rule) Rule {
	parts := strings.Fields(input)

	for i := 2; i < len(parts); i++ {
		switch parts[i] {
		case "--src-type":
			if i+1 < len(parts) {
				ruleItem.Filter.From.AddrType = parts[i+1]
				i++
			}
		case "--dst-type":
			if i+1 < len(parts) {
				ruleItem.Filter.To.AddrType = parts[i+1]
				i++
			}
		}
	}

	return ruleItem
}
