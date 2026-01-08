package handlers

import (
	"context"
	"fmt"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/conntrack"
	"github.com/k0kubun/pp"
)

type ConntrackHandler struct {
}

func NewContrackHandler() *ConntrackHandler {
	return &ConntrackHandler{}
}

func (h *ConntrackHandler) Handle(ctx context.Context) {
	backend := conntrack.NewConntrackBackend()
	list, err := backend.Fetch()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	pp.Println(list)
}
