//go:build SOCKET_DIAG

package socket

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"
	"syscall"

	anetlink "github.com/mdlayher/netlink"
	"github.com/prometheus/procfs"
	"golang.org/x/sys/unix"
)

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

type SocketBackend struct {
}

type SocketDesc struct {
	Port        int
	PID         int
	ListeningIP string
	Comm        string
	CmdLine     string
}

func NewSocketBackend() *SocketBackend {
	return &SocketBackend{}
}

func ipFromBinaryInet(_ [4]uint32) net.IP {
	return net.ParseIP("1.1.1.1")
}

// LIke ss -lptn
func (s *SocketBackend) ListListeners() []SocketDesc {
	c, err := anetlink.Dial(unix.NETLINK_SOCK_DIAG, nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	req := InetDiagReqV2{
		Family:   syscall.AF_INET,
		Protocol: syscall.IPPROTO_TCP,
		States:   TCPF_LISTEN,
	}

	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, &req)

	msgs, err := c.Execute(anetlink.Message{
		Header: anetlink.Header{
			Type:  SOCK_DIAG_BY_FAMILY,
			Flags: anetlink.Request | anetlink.Dump,
		},
		Data: b.Bytes(),
	})
	if err != nil {
		panic(err)
	}

	listSockets := []SocketDesc{}
	for _, m := range msgs {
		var diag InetDiagMsg
		binary.Read(bytes.NewReader(m.Data), binary.LittleEndian, &diag)

		port := binary.BigEndian.Uint16(diag.ID.Sport[:])
		proc, errProcess := findProcessByInode(diag.Inode)
		if errProcess == nil {
			comm, _ := proc.Comm()
			cmdLine, _ := proc.CmdLine()
			listSockets = append(listSockets, SocketDesc{
				Port:        int(port),
				PID:         proc.PID,
				ListeningIP: ipFromBinaryInet(diag.ID.Src).String(),
				CmdLine:     strings.Join(cmdLine, " "),
				Comm:        comm,
			})
		}
	}

	return listSockets
}

func findProcessByInode(inode uint32) (procfs.Proc, error) {
	fs, err := procfs.NewFS("/proc")
	if err != nil {
		return procfs.Proc{}, err
	}

	procs, err := fs.AllProcs()
	if err != nil {
		return procfs.Proc{}, err
	}

	socketTarget := fmt.Sprintf("socket:[%d]", inode)

	for _, p := range procs {
		targets, err := p.FileDescriptorTargets()
		if err != nil {
			return procfs.Proc{}, err
		}

		for _, target := range targets {
			if target == socketTarget {
				if _, err := p.Comm(); err != nil {
					return procfs.Proc{}, err
				}
				return p, nil
			}
		}
	}

	return procfs.Proc{}, errors.New("not found process related to socket inode")
}
