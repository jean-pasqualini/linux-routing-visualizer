package handlers

import (
	"context"
	"fmt"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/routing"
	"github.com/k0kubun/pp"
)

type RoutingHandler struct {
}

func NewRoutingHandler() *RoutingHandler {
	return &RoutingHandler{}
}

func (h *RoutingHandler) Handle(ctx context.Context) {
	fmt.Println("routing")
	rBackend := routing.NewRoutingBackend()
	config, err := rBackend.Fetch()
	if err != nil {
		fmt.Println(err)
	}
	pp.Println(config)
}
