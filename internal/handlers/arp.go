package handlers

import (
	"context"
	"fmt"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/arp"
	"github.com/k0kubun/pp"
)

type ArpHandler struct {
}

func NewArpHandler() *ArpHandler {
	return &ArpHandler{}
}

func (h *ArpHandler) Handle(ctx context.Context) {
	backend := arp.NewArpBackend()
	list, err := backend.Fetch()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	pp.Println(list)
}
