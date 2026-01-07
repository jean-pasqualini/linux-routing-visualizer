package handlers

import (
	"context"
	"fmt"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/sysctl"
	"github.com/k0kubun/pp"
)

type SysctlHandler struct {
}

func NewSysctlHandler() *SysctlHandler {
	return &SysctlHandler{}
}

func (h *SysctlHandler) Handle(ctx context.Context) {
	backend := sysctl.NewSysctlBackend()
	conf, err := backend.Fetch()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	pp.Println(conf)
}
