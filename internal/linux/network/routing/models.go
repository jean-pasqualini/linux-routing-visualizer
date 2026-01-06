package routing

import "net"

type RoutingConfig struct {
	Tables []RoutingTable
	Rules  []WhatRouteTable
}

type RoutingTable struct {
	ID     int
	Routes []RouteDesc
}

type WhatRouteTable struct {
	Priority int
	Src      *net.IPNet
	FwMark   uint32
	Table    RoutingTable
}

type RouteDesc struct {
	Scope      int
	TargetCIDR *net.IPNet
	Device     int
}
