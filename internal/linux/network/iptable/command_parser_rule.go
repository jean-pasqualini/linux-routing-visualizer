package iptable

import (
	"errors"
	"net"
	"strings"

	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
)

func (b *iptableBackend) parseRule(input string) (Rule, error) {
	parts := strings.Fields(input)
	if len(parts) < 3 {
		return Rule{}, errors.New("invalid iptables Rule")
	}
	chainName := parts[1]
	ruleItem := Rule{
		Chain: chainName,
		Raw:   input,
	}

	for i := 2; i < len(parts); i++ {
		switch parts[i] {
		case "-j":
			if i+1 < len(parts) {
				ruleItem.JumpTarget = parts[i+1]
				switch parts[i+1] {
				case "REJECT":
					ruleItem.Action = b.parseActionREJECT(input)
				case "LOG":
					ruleItem.Action = b.parseActionLOG(input)
				case "DNAT", "SNAT":
					ruleItem.Action = b.parseActionNAT(input)
				case "NFLOG":
					ruleItem.Action = b.parseActionNFLOG(input)
				case "TEE":
					ruleItem.Action = b.parseActionTEE(input)
				case "MARK":
					ruleItem.Action = b.parseActionMARK(input)
				case "CONNMARK":
					ruleItem.Action = b.parseActionCONNMARK(input)
				}
				i++
			}
		case "-p":
			if i+1 < len(parts) {
				ruleItem.Filter.Protocol = parts[i+1]
				i++
			}
		case "-s":
			if i+1 < len(parts) {
				_, CIDR, err := net.ParseCIDR(parts[i+1])
				if err != nil {
					state.Dispatch("logger", "!!!!!!!!!!!! PARSING CIDR FAIL : "+err.Error())
				}
				ruleItem.Filter.From.CIDR = CIDR
				i++
			}
		case "-d":
			if i+1 < len(parts) {
				_, CIDR, _ := net.ParseCIDR(parts[i+1])
				ruleItem.Filter.To.CIDR = CIDR
				i++
			}
		case "-i":
			if i+1 < len(parts) {
				negated := false
				if parts[i-1] == "!" {
					negated = true
				}
				ruleItem.Filter.From.Device = Match[string]{Value: parts[i+1], negated: negated}
				i++
			}
		case "-o":
			if i+1 < len(parts) {
				negated := false
				if parts[i-1] == "!" {
					negated = true
				}
				ruleItem.Filter.To.Device = Match[string]{Value: parts[i+1], negated: negated}
				i++
			}
		case "-m":
			if i+1 < len(parts) {
				ruleItem.Modules = append(ruleItem.Modules, parts[i+1])
				switch parts[i+1] {
				case "tcp":
					ruleItem = b.parseRuleModuleTCPUDP(input, ruleItem)
				case "udp":
					ruleItem = b.parseRuleModuleTCPUDP(input, ruleItem)
				case "conntrack":
					ruleItem = b.parseRuleModuleConntrack(input, ruleItem)
				case "addrtype":
					ruleItem = b.parseRuleModuleAddrType(input, ruleItem)
				}
			}
		}
	}

	return ruleItem, nil
}
