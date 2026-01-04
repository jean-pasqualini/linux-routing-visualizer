package handlers

import (
	"context"
	"fmt"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/sniffing"
	"time"
)

type libpcapHandler struct {
	opts LibpcapOptions
}
type LibpcapOptions struct {
	Interface string
	Filter    string
	Snaplen   int32
	Promisc   bool
	Timeout   time.Duration
}

func NewLibpcapHandler(opts LibpcapOptions) *libpcapHandler {
	return &libpcapHandler{opts: opts}
}

func (h *libpcapHandler) Handle(ctx context.Context) {
	iface := h.opts.Interface
	filter := h.opts.Filter
	snaplen := h.opts.Snaplen
	promisc := h.opts.Promisc
	timeout := h.opts.Timeout

	s := sniffing.NewSniffing()
	err := s.Sniff(ctx, iface, filter, snaplen, promisc, timeout)
	if err != nil {
		fmt.Println(err)
	}
}
