package conntrack

import (
	"github.com/florianl/go-conntrack"
)

func convertIPTupple(f *conntrack.IPTuple) IPTupple {
	tracked := IPTupple{}
	if f == nil {
		return tracked
	}
	if f.Src != nil {
		tracked.SrcIP = f.Src.String()
	}
	if f.Dst != nil {
		tracked.DstIP = f.Dst.String()
	}

	if f.Proto != nil {
		if f.Proto.SrcPort != nil {
			tracked.SrcPort = int(*f.Proto.SrcPort)
		}
		if f.Proto.DstPort != nil {
			tracked.DstPort = int(*f.Proto.DstPort)
		}
	}

	return tracked
}

func convert(f conntrack.Con) ConnectionTracked {
	tracked := ConnectionTracked{}
	//pp.Println(f.NatSrc)
	tracked.Translation = convertTranslation(f.Origin, f.Reply)
	tracked.Origin = convertIPTupple(f.Origin)
	tracked.Return = convertIPTupple(f.Reply)
	tracked.Status = convertStatus(f.Status)
	return tracked
}

func convertTranslation(orig *conntrack.IPTuple, reply *conntrack.IPTuple) string {
	if !orig.Src.Equal(*reply.Dst) && !orig.Dst.Equal(*reply.Src) {
		return "SDNAT"
	}
	if !orig.Src.Equal(*reply.Dst) {
		return "SNAT"
	}
	if !orig.Dst.Equal(*reply.Src) {
		return "DNAT"
	}
	return "NONE"
}

func convertStatus(status *uint32) []string {
	if status == nil {
		return nil
	}

	v := *status
	var out []string

	for bit, name := range ctStatusFlags {
		if v&bit != 0 {
			out = append(out, name)
		}
	}

	if len(out) == 0 {
		out = append(out, "NONE")
	}

	return out
}
