package bridge

import (
	"fmt"
	"github.com/vishvananda/netlink"
)

type BridgeBackend struct {
}

func NewBridgeBackend() *BridgeBackend {
	return &BridgeBackend{}
}

type BridgeSpec struct {
	Name  string
	Ports []string
}

func (b *BridgeBackend) convert(bridge *netlink.Bridge, linkList []netlink.Link) BridgeSpec {
	spec := BridgeSpec{
		Name: bridge.Name,
	}

	for _, link := range linkList {
		if link.Attrs().MasterIndex == bridge.Attrs().Index {
			spec.Ports = append(spec.Ports, fmt.Sprintf("%s(type=%s)", link.Attrs().Name, link.Type()))
		}
	}

	return spec
}

func (b *BridgeBackend) Fetch() ([]BridgeSpec, error) {
	list := []BridgeSpec{}

	links, err := netlink.LinkList()
	if err != nil {
		return list, err
	}
	for _, link := range links {
		if bridge, ok := link.(*netlink.Bridge); ok {
			list = append(list, b.convert(bridge, links))
		}
	}

	return list, nil
}
