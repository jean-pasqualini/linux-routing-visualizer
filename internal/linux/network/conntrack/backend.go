package conntrack

import "github.com/florianl/go-conntrack"

type ConntrackBackend struct {
}

func NewConntrackBackend() *ConntrackBackend {
	return &ConntrackBackend{}
}

type ConnectionTracked struct {
	SrcIP   string
	DstIP   string
	SrcPort int
	DstPort int
	Status  int
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
		list = append(list, ConnectionTracked{
			SrcIP:   f.Origin.Src.String(),
			DstIP:   f.Origin.Dst.String(),
			SrcPort: int(*f.Origin.Proto.SrcPort),
			DstPort: int(*f.Origin.Proto.DstPort),
			Status:  int(*f.Status),
		})
	}

	return list, nil
}
