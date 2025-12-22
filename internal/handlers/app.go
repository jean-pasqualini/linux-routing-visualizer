package handlers

import (
	"context"
	"fmt"

	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/iptable"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/logging"
)

type AppHandler struct {
	iptReader *iptable.IPtableReader
}

func NewAppHandler() *AppHandler {
	return &AppHandler{
		iptReader: iptable.NewIPtableReader(),
	}
}

func (h *AppHandler) Handle(context context.Context) {
	logger := logging.FromContext(context)
	logger.Debug("Handling request")
	fmt.Println("Hello World")

	h.iptReader.Read(context)
}
