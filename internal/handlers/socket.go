package handlers

import (
	"context"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/socket"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/logging"
	"github.com/k0kubun/pp"
)

type SocketHandler struct {
}

func NewSocketHandler() *SocketHandler {
	return &SocketHandler{}
}

func (h *SocketHandler) Handle(context context.Context) {
	logger := logging.FromContext(context)
	logger.Info("socekt handler")

	sBackend := socket.NewSocketBackend()
	socketList := sBackend.ListListeners()
	pp.Println(socketList)
}
