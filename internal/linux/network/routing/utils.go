package routing

import (
	"fmt"
	"net"
	"syscall"

	"github.com/vishvananda/netlink"
)

func IsLocalRoute(ip net.IP) bool {
	routes, _ := netlink.RouteGet(ip)
	return routes[0].Type == syscall.RTN_LOCAL
}

func GetIPFromNetDevice(name string) (net.IP, error) {
	// 2) interface de sortie
	link, err := netlink.LinkByName(name)
	if err != nil {
		return nil, err
	}

	// 3) adresse source associée
	addrs, err := netlink.AddrList(link, netlink.FAMILY_V4)
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		if addr.IP != nil && !addr.IP.IsLoopback() {
			return addr.IP, nil
		}
	}

	return nil, fmt.Errorf("aucune IP trouvée")
}

func RervesePathCheck(srcIP net.IP, inIfName string, mode int) bool {
	if mode == 0 {
		return true
	}

	inLink, err := netlink.LinkByName(inIfName)
	if err != nil {
		return false
	}

	opts := &netlink.RouteGetOptions{
		SrcAddr:  srcIP,
		IifIndex: inLink.Attrs().Index,
		FIBMatch: true,
	}

	routes, err := netlink.RouteGetWithOptions(srcIP, opts)
	if err != nil {
		return false
	}
	if len(routes) < 1 {
		return false
	}
	rev := routes[0]

	// // loose: only require route existence
	if mode == 2 {
		return true
	}

	// mode == 1 strict: require oif == iif
	if rev.LinkIndex == inLink.Attrs().Index {
		return true
	}
	return false
}
