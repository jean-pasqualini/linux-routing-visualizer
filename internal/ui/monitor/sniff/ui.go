package sniff

import (
	"fmt"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/sniffing"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/rivo/tview"
)

type SniffView struct {
	*tview.TextView
}

func NewSniffView() *SniffView {
	textView := tview.NewTextView()
	textView.SetBorder(true).SetTitle("Sniffing")
	v := SniffView{textView}

	state.Subscribe("monitor:packet", v.OnPacket)

	return &v
}

func (v *SniffView) OnPacket(name string, event any) {
	if packet, ok := event.(sniffing.Packet); ok {
		fmt.Fprintf(v, "TCP FROM %s to %s\n", packet.Source, packet.Target)
	}
}
