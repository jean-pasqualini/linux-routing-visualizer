package handlers

import (
	"context"
	"fmt"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/bridge"
	"github.com/k0kubun/pp"
)

type BridgeHandler struct {
}

func NewBridgeHandler() *BridgeHandler {
	return &BridgeHandler{}
}

func (h *BridgeHandler) Handle(ctx context.Context) {
	backend := bridge.NewBridgeBackend()
	list, err := backend.Fetch()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	pp.Println(list)
}
