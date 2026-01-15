package iptable

type ActionReject struct {
	RejectWith string
}

type ActionLog struct {
	LogPrefix string
	LogLevel  int
}

type ActionNat struct {
	IP string
}

type ActionMark struct {
	Mark string
}

type ActionConnMark struct {
	SaveMark bool
}

type ActionNFLOG struct {
	Group int
}

type ActionTEE struct {
	Gateway string
}
