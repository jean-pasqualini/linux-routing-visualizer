package simulator

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/iptable"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/simulator"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/vishvananda/netlink"
	"net"
	"syscall"
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
			Rules:   []simulator.SimulatorResultRuleEvent{},
		},
	}
}

func (s *IptableSimulator) Match(query simulator.FormEvent, rule iptable.Rule) bool {
	if rule.Filter.To.Port != "" && rule.Filter.To.Port != query.TargetPort {
		return false
	}
	if rule.Filter.From.CIDR != nil && false {
		return false
	}
	if rule.Filter.To.CIDR != nil && rule.Filter.To.CIDR.Contains(net.ParseIP(query.TargetIP)) {
		return false
	}
	return true
}

func (s *IptableSimulator) enterChain(tableName string, chainName string) {
	if chain, ok := s.tables[tableName].Chains[chainName]; ok {
		state.Dispatch("logger", "\t -> enter chain "+chainName)
		for _, rule := range chain.Rules {
			ruleMatching := simulator.SimulatorResultRuleEvent{
				Raw:     rule.Raw,
				Matched: s.Match(s.query, rule),
			}
			state.Dispatch("logger", "\t\t -> should add rule "+ruleMatching.Raw)
			s.result.Rules = append(s.result.Rules, ruleMatching)
			if ruleMatching.Matched {
				s.enterChain(tableName, rule.JumpTarget)
			}
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
	targetIP := net.ParseIP(s.query.TargetIP)
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
