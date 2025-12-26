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
	Chain      string
	JumpTarget string
	Filter     ruleFilter
}

type ruleFilter struct {
	Protocol string
	To       ruleFilterFromTo
}

type ruleFilterFromTo struct {
	Port   string
	Device string
}

type counter struct {
	Packets uint64
	Bytes   uint64
}
