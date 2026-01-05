package link

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/k0kubun/pp"
	anetlink "github.com/mdlayher/netlink"
	"github.com/prometheus/procfs"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"

	"syscall"
)

type LinkBackend struct {
	links []netlink.Link
}

func NewLinkBackend() *LinkBackend {
	links, _ := netlink.LinkList()
	return &LinkBackend{
		links: links,
	}
}

func (b *LinkBackend) GetInterfacesNames() []string {
	names := []string{}
	for _, link := range b.links {
		names = append(names, link.Attrs().Name)
	}

	return names
}

func Fetch() {
	links, _ := netlink.LinkList()

	pp.Println(links)

	for _, link := range links {
		fmt.Printf("Interface %s\n", link.Attrs().Name)
	}
}

// LIke ss -lptn
func ListListeners() {
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

	for _, m := range msgs {
		var diag InetDiagMsg
		binary.Read(bytes.NewReader(m.Data), binary.LittleEndian, &diag)

		port := binary.BigEndian.Uint16(diag.ID.Sport[:])
		fmt.Println("PORT:", port, "INODE:", diag.Inode)
		pid, program, errProcess := findProcessByInode(diag.Inode)
		pp.Println(pid, program, errProcess)
		if errProcess != nil {
			fmt.Printf("found => %d %s", pid, program)
		} else {
			pp.Println(errProcess)
		}
	}
}

func findProcessByInode(inode uint32) (int, string, error) {
	fs, err := procfs.NewFS("/proc")
	if err != nil {
		return 0, "", err
	}

	procs, err := fs.AllProcs()
	if err != nil {
		return 0, "", err
	}

	socketTarget := fmt.Sprintf("socket:[%d]", inode)

	for _, p := range procs {
		targets, err := p.FileDescriptorTargets()
		if err != nil {
			return 0, "", err
		}

		for _, target := range targets {
			if target == socketTarget {
				comm, err := p.Comm()
				if err != nil {
					return 0, "", err
				}
				return p.PID, comm, nil
			}
		}
	}

	return 0, "", nil
}
