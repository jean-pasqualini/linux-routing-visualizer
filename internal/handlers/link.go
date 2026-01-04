package handlers

import (
	"context"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/link"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/logging"
)

type LinkHandler struct {
}

func NewLinkHandler() *LinkHandler {
	return &LinkHandler{}
}

func (h *LinkHandler) Handle(context context.Context) {
	logger := logging.FromContext(context)
	logger.Info("link handler")

	link.Fetch()
}
