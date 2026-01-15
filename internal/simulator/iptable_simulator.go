package simulator

import (
	"net"
	"slices"
	"strings"

	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/iptable"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/routing"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
)

type SimulatorQuery struct {
	Source           SimulatorQueryTarget
	Target           SimulatorQueryTarget
	IncludeUnmatched bool
	State            string
	Protocol         string
}

type SimulatorQueryTarget struct {
	Device string
	IP     string
	Port   string
}

type Simulator struct {
	query  SimulatorQuery
	tables map[string]iptable.Table
	result SimulatorResult
}

type SimulatorMismatch struct {
	Reason string
	Raw    string
}

type SimulatorResultRuleEvent struct {
	Raw        string
	JumpChain  *SimulatorNetfilterChainEvent
	Action     string
	Matched    bool
	Mismatches []SimulatorMismatch
}

type SimulatorNetfilterNatEvent struct {
	OldIP string
	NewIP string
}

type SimulatorNetfilterChainEvent struct {
	Name     string
	Decision string
	Events   []SimulatorEvent
}

type SimulatorNetrouting struct {
	RouteType string
}

type SimulatorIncomingInterface struct {
	Interface string
}

type SimulatorEvent interface {
}

type SimulatorResult struct {
	Events []SimulatorEvent
}

func NewSimulator(query SimulatorQuery, tables map[string]iptable.Table) *Simulator {
	return &Simulator{
		query:  query,
		tables: tables,
		result: SimulatorResult{
			Events: make([]SimulatorEvent, 0),
		},
	}
}

func (s *Simulator) Match(query SimulatorQuery, rule iptable.Rule) (bool, []SimulatorMismatch) {
	reasons := make([]SimulatorMismatch, 0)
	// AddrType
	if rule.Filter.From.AddrType == "LOCAL" && !routing.IsLocalRoute(net.ParseIP(query.Source.IP)) {
		reasons = append(reasons, SimulatorMismatch{"source addrtype", "--src-type"})
	}
	if rule.Filter.To.AddrType == "LOCAL" && !routing.IsLocalRoute(net.ParseIP(query.Target.IP)) {
		reasons = append(reasons, SimulatorMismatch{"target addrtype", "--dst-type"})
	}
	// Protocl
	if rule.Filter.Protocol != "" && strings.ToUpper(rule.Filter.Protocol) != query.Protocol {
		reasons = append(reasons, SimulatorMismatch{"protocol mismatch ", "-p"})
	}
	// Device
	if !rule.Filter.From.Device.Match(query.Source.Device) {
		reasons = append(reasons, SimulatorMismatch{"source device mismatch", "-i"})
	}
	if !rule.Filter.To.Device.Match(query.Target.Device) {
		reasons = append(reasons, SimulatorMismatch{"target device mismatch", "-o"})
	}
	// Port
	if rule.Filter.From.Port != "" && rule.Filter.From.Port != query.Source.Port {
		reasons = append(reasons, SimulatorMismatch{"source port mismatch", "--sport"})
	}
	if rule.Filter.To.Port != "" && rule.Filter.To.Port != query.Target.Port {
		reasons = append(reasons, SimulatorMismatch{"target port mismatch", "--dport"})
	}
	// Cidr
	if rule.Filter.From.CIDR != nil && !rule.Filter.From.CIDR.Contains(net.ParseIP(query.Source.IP)) {
		reasons = append(reasons, SimulatorMismatch{"source cidr mismatch", "-s"})
	}
	if rule.Filter.To.CIDR != nil && !rule.Filter.To.CIDR.Contains(net.ParseIP(query.Target.IP)) {
		reasons = append(reasons, SimulatorMismatch{"target cidr mismatch", "-t"})
	}
	if len(rule.Filter.ConnectionState) > 0 && !slices.Contains(rule.Filter.ConnectionState, query.State) {
		reasons = append(reasons, SimulatorMismatch{"conection state mismatch", "--cstate"})
	}
	return len(reasons) == 0, reasons
}

func (s *Simulator) isJumpAction(target string) bool {
	return slices.Contains([]string{"ACCEPT", "DROP", "REJECT", "RETURN", "TRACE", "SNAT", "DNAT", "MASQUERADE"}, target)
}

func (s *Simulator) enterChain(chainEvent *SimulatorNetfilterChainEvent, tableName string, chainName string) {
	if chain, ok := s.tables[tableName].Chains[chainName]; ok {
		state.Dispatch("simulator:log", "\t -> enter chain "+chainName)
		for _, rule := range chain.Rules {

			matchingResult, mismatches := s.Match(s.query, rule)
			ruleMatching := SimulatorResultRuleEvent{
				Raw:        rule.Raw,
				Matched:    matchingResult,
				Mismatches: mismatches,
			}
			state.Dispatch("simulator:log", "\t\t -> should add rule "+ruleMatching.Raw)
			if !matchingResult {
				state.Dispatch("simulator:log", "\t\t -> matching fail")
				for _, mismatch := range mismatches {
					state.Dispatch("simulator:log", "\t\t\t -> filter: "+mismatch.Reason)
				}
			}

			if ruleMatching.Matched || s.query.IncludeUnmatched {
				chainEvent.Events = append(chainEvent.Events, ruleMatching)
			}

			if ruleMatching.Matched {
				if rule.JumpTarget == "ACCEPT" || rule.JumpTarget == "DROP" || rule.JumpTarget == "REJECT" || rule.JumpTarget == "RETURN" {
					state.Dispatch("simulator:log", "\t\t\t -> final decision "+rule.JumpTarget+" for that chain ")
					chainEvent.Decision = rule.JumpTarget
					return
				}
				if rule.JumpTarget == "TRACE" {
					ruleMatching.Action = "TRACE"
					state.Dispatch("simulator:log", "\t\t\t -> tracing enable")
				} else if rule.JumpTarget == "DNAT" {
					ruleMatching.Action = "DNAT"
					state.Dispatch("simulator:log", "\t\t\t -> dnat")
				} else if rule.JumpTarget == "MASQUERADE" {
					ip, err := routing.GetIPFromNetDevice(s.query.Target.Device)
					if err == nil {
						chainEvent.Events = append(chainEvent.Events, SimulatorNetfilterNatEvent{
							OldIP: s.query.Source.IP,
							NewIP: ip.String(),
						})
					}
					ruleMatching.Action = "MASQUERADE"
					state.Dispatch("simulator:log", "\t\t\t -> masquerade")
				} else {
					ruleMatching.Action = "JUMP"
					jumpChain := &SimulatorNetfilterChainEvent{Name: rule.JumpTarget, Decision: "NONE"}
					s.enterChain(jumpChain, tableName, rule.JumpTarget)
					chainEvent.Events = append(chainEvent.Events, *jumpChain)
				}
			}
		}
		if chain.Policy != "-" {
			chainEvent.Decision = chain.Policy
			state.Dispatch("simulator:log", "\t\t\t -> final decision "+chain.Policy+" for that chain ")
		} else {
			chainEvent.Decision = "NONE"
		}
	} else {
		state.Dispatch("simulator:log", "\t -> not found chain "+chainName)
	}
}

func (s *Simulator) enterTable(chainEvent *SimulatorNetfilterChainEvent, tableName string) {
	if _, ok := s.tables[tableName]; ok {
		state.Dispatch("simulator:log", "-> enter table "+tableName)
		s.enterChain(chainEvent, tableName, chainEvent.Name)
	} else {
		state.Dispatch("simulator:log", "-> not found table "+tableName)
	}
}

func (s *Simulator) walkBuiltinChain(chainEvent *SimulatorNetfilterChainEvent) {
	for _, tableName := range iptable.TablesList {
		s.enterTable(chainEvent, string(tableName))
	}
}

func (s *Simulator) Simulate() SimulatorResult {
	sourceIP := net.ParseIP(s.query.Target.IP)
	if sourceIP == nil {
		state.Dispatch("simulator:log", "invalid source ip")
		return SimulatorResult{}
	}
	targetIP := net.ParseIP(s.query.Target.IP)
	if targetIP == nil {
		state.Dispatch("simulator:log", "invalid target ip")
		return SimulatorResult{}
	}
	s.result.Events = append(s.result.Events, SimulatorIncomingInterface{
		Interface: s.query.Source.Device,
	})

	// We consider for now that it is an incoming packet
	// Chain PREPROUTING
	preroutingEvent := SimulatorNetfilterChainEvent{Name: "PREROUTING"}
	s.walkBuiltinChain(&preroutingEvent)
	s.result.Events = append(s.result.Events, preroutingEvent)

	if routing.IsLocalRoute(targetIP) {
		s.result.Events = append(s.result.Events, SimulatorNetrouting{
			RouteType: "local",
		})
		// Chain INPUT
		inputEvent := SimulatorNetfilterChainEvent{Name: "INPUT"}
		s.walkBuiltinChain(&inputEvent)
		s.result.Events = append(s.result.Events, inputEvent)
	} else {
		s.result.Events = append(s.result.Events, SimulatorNetrouting{
			RouteType: "nolocal",
		})
		// Chain FORWARD
		forwardEvent := SimulatorNetfilterChainEvent{Name: "FORWARD"}
		s.walkBuiltinChain(&forwardEvent)
		s.result.Events = append(s.result.Events, forwardEvent)
		postRoutingEvent := SimulatorNetfilterChainEvent{Name: "POSTROUTING"}
		s.walkBuiltinChain(&postRoutingEvent)
		s.result.Events = append(s.result.Events, postRoutingEvent)
	}

	return s.result
}
