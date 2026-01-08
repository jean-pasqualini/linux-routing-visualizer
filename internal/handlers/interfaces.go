package handlers

import (
	"context"
	"fmt"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/netdevice"
	"github.com/k0kubun/pp"
)

type InterfacesHandler struct {
}

func NewInterfacesHandler() *InterfacesHandler {
	return &InterfacesHandler{}
}

func (h *InterfacesHandler) Handle(ctx context.Context) {
	backend := netdevice.NewInterfacesBackend()
	list, err := backend.Fetch()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	pp.Println(list)
}
