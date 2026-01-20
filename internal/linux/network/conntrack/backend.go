package conntrack

import (
	"context"
	"github.com/florianl/go-conntrack"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
)

type ConntrackBackend struct {
}

func NewConntrackBackend() *ConntrackBackend {
	return &ConntrackBackend{}
}

func (b *ConntrackBackend) Fetch() ([]ConnectionTracked, error) {
	list := []ConnectionTracked{}
	nfct, err := conntrack.Open(&conntrack.Config{})
	if err != nil {
		return list, err
	}
	defer nfct.Close()

	flows, err := nfct.Dump(conntrack.Conntrack, conntrack.IPv4)
	if err != nil {
		return list, err
	}

	for _, f := range flows {
		list = append(list, convert(f))
	}

	return list, nil
}

func (b *ConntrackBackend) FetchLive(ctx context.Context) (chan ConnectionTracked, error) {
	out := make(chan ConnectionTracked, 4096)

	go func() {
		nfct, err := conntrack.Open(&conntrack.Config{})
		if err != nil {
			close(out)
			state.Dispatch("app:log", "oepn conntrack live: "+err.Error())
			return
		}
		defer nfct.Close()
		if err := nfct.Register(
			ctx,
			conntrack.Conntrack,
			conntrack.NetlinkCtNew, func(c conntrack.Con) int {
				out <- convert(c)
				return 0
			}); err != nil {
			close(out)
			state.Dispatch("app:log", "register conntrack live: "+err.Error())
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, nil
}
