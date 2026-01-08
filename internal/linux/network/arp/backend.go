package arp

import "github.com/vishvananda/netlink"

type ArpBackend struct {
}

func NewArpBackend() *ArpBackend {
	return &ArpBackend{}
}

type ArpEntry struct {
	IP    string
	HAddr string
}

func (b *ArpBackend) Fetch() ([]ArpEntry, error) {
	list := []ArpEntry{}
	neights, err := netlink.NeighList(0, netlink.FAMILY_V4)
	if err != nil {
		return list, err
	}
	for _, n := range neights {
		list = append(list, ArpEntry{
			IP:    n.IP.String(),
			HAddr: n.HardwareAddr.String(),
		})
	}

	return list, nil
}
