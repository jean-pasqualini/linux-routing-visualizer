package link

import (
	"fmt"

	"github.com/k0kubun/pp"
	"github.com/vishvananda/netlink"
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
