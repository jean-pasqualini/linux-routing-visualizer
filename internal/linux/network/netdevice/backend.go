package netdevice

import (
	"github.com/vishvananda/netlink"
	"net"
)

type InterfacesBackend struct {
}

func NewInterfacesBackend() *InterfacesBackend {
	return &InterfacesBackend{}
}

type NetDeviceSpec struct {
	Name     string
	HAddr    string
	Type     string
	AddrList []string
	Up       bool
}

func (b *InterfacesBackend) Fetch() ([]NetDeviceSpec, error) {
	list := []NetDeviceSpec{}
	links, err := netlink.LinkList()
	if err != nil {
		return list, err
	}
	for _, link := range links {
		addrList := []string{}
		addrs, err := netlink.AddrList(link, netlink.FAMILY_V4)
		if err == nil {
			for _, addr := range addrs {
				addrList = append(addrList, addr.IP.String())
			}
		}
		list = append(list, NetDeviceSpec{
			Name:     link.Attrs().Name,
			HAddr:    link.Attrs().HardwareAddr.String(),
			Type:     link.Type(),
			AddrList: addrList,
			Up:       link.Attrs().Flags&net.FlagUp != 0,
		})
	}

	return list, nil
}
