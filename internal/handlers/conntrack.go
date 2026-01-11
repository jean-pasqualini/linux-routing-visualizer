package handlers

import (
	"context"
	"fmt"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/conntrack"
)

type ConntrackHandler struct {
}

func NewContrackHandler() *ConntrackHandler {
	return &ConntrackHandler{}
}

func (h *ConntrackHandler) Handle(ctx context.Context) {
	backend := conntrack.NewConntrackBackend()
	_, err := backend.Fetch()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//pp.Println(list)
}
