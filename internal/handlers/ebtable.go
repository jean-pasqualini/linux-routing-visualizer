package handlers

import (
	"context"
	"fmt"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/ebtable"
	"github.com/k0kubun/pp"
)

type EbtableHandler struct {
}

func NewEbtableHandler() *EbtableHandler {
	return &EbtableHandler{}
}

func (h *EbtableHandler) Handle(ctx context.Context) {
	backend := ebtable.NewEbtableBackend()
	table, err := backend.Fetch()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	pp.Println(table)
}
