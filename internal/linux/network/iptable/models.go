package iptable

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
	Modules    []string
	Chain      string
	JumpTarget string
	Filter     ruleFilter
}

type ruleFilter struct {
	Protocol        string
	ConnectionState string
	To              ruleFilterFromTo
	From            ruleFilterFromTo
}

type ruleFilterFromTo struct {
	Port     string
	AddrType string
	Device   string
}

type counter struct {
	Packets uint64
	Bytes   uint64
}

type protocolType string

const (
	tcp    protocolType = "tcp"
	udp    protocolType = "udp"
	icmp   protocolType = "icmp"
	icmpv6 protocolType = "icmpv6"
	ipip   protocolType = "ipip"
)

type TableType string

const (
	raw      TableType = "raw"      // Before conntrack
	mangle   TableType = "mangle"   // To modify packet (TTL, marks, Qos)
	nat      TableType = "nat"      // To SNAT/DNAT
	filter   TableType = "filter"   // To filter out
	security TableType = "security" // SELinux
)

type ChainType string

const (
	PREROUTING  ChainType = "PREROUTING"
	INPUT       ChainType = "INPUT"
	FORWARD     ChainType = "FORWARD"
	OUTPUT      ChainType = "OUTPUT"
	POSTROUTING ChainType = "POSTROUTING"
)

type TargetAction string

const (
	ACCEPT     TargetAction = "ACCEPT"
	DROP       TargetAction = "DROP"
	REJECT     TargetAction = "REJECT"
	LOG        TargetAction = "LOG"
	DNAT       TargetAction = "DNAT"
	SNAT       TargetAction = "SNAT"
	MASQUERADE TargetAction = "MASQUERADE"
	RETURN     TargetAction = "RETURN"
)

var actionList = [...]TargetAction{ACCEPT, DROP, REJECT, LOG, DNAT, SNAT, MASQUERADE, RETURN}

var rawBuiltinChains = [...]ChainType{PREROUTING, OUTPUT}
var mangleBuiltinChains = [...]ChainType{PREROUTING, INPUT, FORWARD, OUTPUT, POSTROUTING}
var natBuiltinChains = [...]ChainType{PREROUTING, INPUT, FORWARD, OUTPUT}
var filterBuiltinChains = [...]ChainType{INPUT, FORWARD, OUTPUT}
var securityBuiltinChains = [...]ChainType{INPUT, FORWARD, OUTPUT}

var InboundChaining = [...]ChainType{PREROUTING, INPUT}
var OutboundChaining = [...]ChainType{OUTPUT, POSTROUTING}
var ForwardChaining = [...]ChainType{PREROUTING, FORWARD, POSTROUTING}

var TablesList = [...]TableType{raw, mangle, nat, filter, security}
