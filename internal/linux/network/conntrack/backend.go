package conntrack

import "github.com/florianl/go-conntrack"

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
