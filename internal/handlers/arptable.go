package handlers

import (
	"context"
	"fmt"

	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/arptable"
	"github.com/k0kubun/pp"
)

type ArptableHandler struct {
}

func NewArptableHandler() *ArptableHandler {
	return &ArptableHandler{}
}

func (h *ArptableHandler) Handle(ctx context.Context) {
	backend := arptable.NewArpTableBackend()
	list, err := backend.Fetch()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	pp.Println(list)
}
