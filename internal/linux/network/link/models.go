package link

const (
	NETLINK_SOCK_DIAG   = 4
	SOCK_DIAG_BY_FAMILY = 20
	TCPF_LISTEN         = 1 << 10
)

type InetDiagReqV2 struct {
	Family   uint8
	Protocol uint8
	Ext      uint8
	Pad      uint8
	States   uint32
	ID       InetDiagSockID
}

type InetDiagSockID struct {
	Sport  [2]byte
	Dport  uint16
	Src    [4]uint32
	Dst    [4]uint32
	If     uint32
	Cookie [2]uint32
}

type InetDiagMsg struct {
	Family  uint8
	State   uint8
	Timer   uint8
	Retrans uint8
	ID      InetDiagSockID
	Expires uint32
	Rqueue  uint32
	Wqueue  uint32
	UID     uint32
	Inode   uint32
}
