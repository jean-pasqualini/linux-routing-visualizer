package handlers

import (
	"context"
	"fmt"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/sniffing"
	"sync"
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
	ch, _ := s.Sniff(ctx, iface, filter, snaplen, promisc, timeout)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for pkt := range ch {
			fmt.Printf("IPV4 %s -> %s \n", pkt.Source, pkt.Target)
		}
	}()
	wg.Wait()
}
