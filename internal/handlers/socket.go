package handlers

import (
	"context"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/link"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/logging"
)

type SocketHandler struct {
}

func NewSocketHandler() *SocketHandler {
	return &SocketHandler{}
}

func (h *SocketHandler) Handle(context context.Context) {
	logger := logging.FromContext(context)
	logger.Info("socekt handler")

	link.ListListeners()
}
