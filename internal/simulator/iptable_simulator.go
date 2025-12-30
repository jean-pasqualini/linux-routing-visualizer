package simulator

import (
	"net"
	"syscall"

	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/iptable"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/simulator"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/vishvananda/netlink"
)

type IptableSimulator struct {
	query  simulator.FormEvent
	tables map[string]iptable.Table
	result simulator.SimulatorResultEvent
}

func NewIptableSimulator(query simulator.FormEvent, tables map[string]iptable.Table) *IptableSimulator {
	return &IptableSimulator{
		query:  query,
		tables: tables,
		result: simulator.SimulatorResultEvent{
			Request: query,
			Rules:   make([]simulator.SimulatorResultRuleEvent, 0),
			Chains:  make([]simulator.SimulatorResultChainEvent, 0),
		},
	}
}

func (s *IptableSimulator) Match(query simulator.FormEvent, rule iptable.Rule) (bool, []string) {
	reasons := make([]string, 0)
	// Protocl
	if rule.Filter.Protocol != "" && rule.Filter.Protocol != query.Protocol {
		state.Dispatch("logger", "\t\t\t -> filter: protocol mismatch")
		reasons = append(reasons, "protocol mismatch")
	}
	// Device
	if !rule.Filter.From.Device.Match(query.Source.Device) {
		state.Dispatch("logger", "\t\t\t -> filter: source device mismatch")
		reasons = append(reasons, "source device mismatch")
	}
	if !rule.Filter.To.Device.Match(query.Target.Device) {
		state.Dispatch("logger", "\t\t\t -> filter: target device mismatch")
		reasons = append(reasons, "target device mismatch")
	}
	// Port
	if rule.Filter.From.Port != "" && rule.Filter.From.Port != query.Source.Port {
		state.Dispatch("logger", "\t\t\t -> filter: source port mismatch")
		reasons = append(reasons, "source port mismatch")
	}
	if rule.Filter.To.Port != "" && rule.Filter.To.Port != query.Target.Port {
		state.Dispatch("logger", "\t\t\t -> filter: target port mismatch")
		reasons = append(reasons, "target port mismatch")
	}
	// Cidr
	if rule.Filter.From.CIDR != nil && !rule.Filter.From.CIDR.Contains(net.ParseIP(query.Source.IP)) {
		state.Dispatch("logger", "\t\t\t -> filter: source cidr mismatch")
		reasons = append(reasons, "source cidr mismatch")
	}
	if rule.Filter.To.CIDR != nil && !rule.Filter.To.CIDR.Contains(net.ParseIP(query.Target.IP)) {
		state.Dispatch("logger", "\t\t\t -> filter: target cidr mismatch")
		reasons = append(reasons, "target cidr mismatch")
	}
	return len(reasons) == 0, reasons
}

func (s *IptableSimulator) enterChain(tableName string, chainName string) {
	if chain, ok := s.tables[tableName].Chains[chainName]; ok {
		state.Dispatch("logger", "\t -> enter chain "+chainName)
		for _, rule := range chain.Rules {

			matchingResult, reasons := s.Match(s.query, rule)
			ruleMatching := simulator.SimulatorResultRuleEvent{
				Raw:     rule.Raw,
				Matched: matchingResult,
			}
			state.Dispatch("logger", "\t\t -> should add rule "+ruleMatching.Raw)
			if !matchingResult {
				state.Dispatch("logger", "\t\t -> matching fail")
				for _, reason := range reasons {
					state.Dispatch("logger", "\t\t\t -> filter: "+reason)
				}
			}

			s.result.Rules = append(s.result.Rules, ruleMatching)
			if ruleMatching.Matched {
				if rule.JumpTarget == "ACCEPT" || rule.JumpTarget == "DROP" || rule.JumpTarget == "REJECT" || rule.JumpTarget == "RETURN" {
					s.result.Chains = append(s.result.Chains, simulator.SimulatorResultChainEvent{
						Name:     chainName,
						Decision: rule.JumpTarget,
					})
					state.Dispatch("logger", "\t\t\t -> final decision "+rule.JumpTarget+" for that chain ")
					return
				}
				if rule.JumpTarget == "TRACE" {
					state.Dispatch("logger", "\t\t\t -> tracing enable")
				}
				if rule.JumpTarget == "DNAT" {
					state.Dispatch("logger", "\t\t\t -> dnat")
				}
				if rule.JumpTarget == "MASQUERADE" {
					state.Dispatch("logger", "\t\t\t -> masquerade")
				}
				s.enterChain(tableName, rule.JumpTarget)
			}
		}
		if chain.Policy != "-" {
			s.result.Chains = append(s.result.Chains, simulator.SimulatorResultChainEvent{
				Name:     chainName,
				Decision: chain.Policy,
			})
			state.Dispatch("logger", "\t\t\t -> final decision "+chain.Policy+" for that chain ")
		}
	} else {
		state.Dispatch("logger", "\t -> not found chain "+chainName)
	}
}

func (s *IptableSimulator) isLocalRoute(ip net.IP) bool {
	routes, _ := netlink.RouteGet(ip)
	return routes[0].Type == syscall.RTN_LOCAL
}

func (s *IptableSimulator) enterTable(tableName string, chainName string) {
	if _, ok := s.tables[tableName]; ok {
		state.Dispatch("logger", "-> enter table "+tableName)
		s.enterChain(tableName, chainName)
	} else {
		state.Dispatch("logger", "-> not found table "+tableName)
	}
}

func (s *IptableSimulator) Simulate() simulator.SimulatorResultEvent {
	sourceIP := net.ParseIP(s.query.Target.IP)
	if sourceIP == nil {
		state.Dispatch("logger", "invalid source ip")
		return simulator.SimulatorResultEvent{}
	}
	targetIP := net.ParseIP(s.query.Target.IP)
	if targetIP == nil {
		state.Dispatch("logger", "invalid target ip")
		return simulator.SimulatorResultEvent{}
	}
	// We consider for now that it is an incoming packet
	// Chain PREPROUTING
	for _, tableName := range iptable.TablesList {
		s.enterTable(string(tableName), "PREROUTING")
	}
	if s.isLocalRoute(targetIP) {
		// Chain INPUT
		for _, tableName := range iptable.TablesList {
			s.enterTable(string(tableName), "INPUT")
		}
	} else {
		// Chain FORWARD
		for _, tableName := range iptable.TablesList {
			s.enterTable(string(tableName), "FORWARD")
		}
		for _, tableName := range iptable.TablesList {
			s.enterTable(string(tableName), "POSTROUTING")
		}
	}

	return s.result
}
