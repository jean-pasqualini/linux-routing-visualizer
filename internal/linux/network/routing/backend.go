package routing

import "github.com/vishvananda/netlink"

type RoutingBackend struct {
}

func NewRoutingBackend() *RoutingBackend {
	return &RoutingBackend{}
}

func (r *RoutingBackend) Fetch() (RoutingConfig, error) {
	config := RoutingConfig{}
	rules, err := netlink.RuleList(netlink.FAMILY_V4)
	if err != nil {
		return config, err
	}
	for _, rule := range rules {
		table := RoutingTable{
			ID: rule.Table,
		}
		wrt := WhatRouteTable{
			Priority: rule.Priority,
			Src:      rule.Src,
			FwMark:   rule.Mark,
			Table:    table,
		}
		routeList := []RouteDesc{}
		routes, err := netlink.RouteListFiltered(netlink.FAMILY_V4, &netlink.Route{Table: rule.Table}, netlink.RT_FILTER_TABLE)
		if err == nil {
			for _, route := range routes {
				routeList = append(routeList, RouteDesc{
					Scope:      int(route.Scope),
					Device:     route.LinkIndex,
					TargetCIDR: route.Dst,
				})
			}
		}
		wrt.Table.Routes = routeList
		config.Rules = append(config.Rules, wrt)
		//config.Tables = append(config.Tables, table)
	}
	return config, nil
}
