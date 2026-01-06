package handlers

import (
	"context"
	"fmt"

	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/ipvs"
	"github.com/k0kubun/pp"
)

type IPVSHandler struct {
}

func NewIPVSHandler() *IPVSHandler {
	return &IPVSHandler{}
}

func (h *IPVSHandler) Handle(ctx context.Context) {
	ipvsBackend := ipvs.NewIPVSBackend()
	raw, services := ipvsBackend.Fetch()
	fmt.Println(raw)
	pp.Print(services)
}
