package monitor

import (
	"fmt"

	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/sniffing"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/rivo/tview"
)

type MonitorView struct {
	*tview.TextView
	started bool
}

func NewMonitorView() *MonitorView {
	mview := &MonitorView{
		TextView: tview.NewTextView(),
	}

	state.Subscribe("monitor:packet", mview.OnPacket)

	return mview
}

func (v *MonitorView) OnTabShow() {
	if v.started {
		state.Dispatch("app:log", "already started")
		return
	}
	v.started = true
	state.Dispatch("monitor:start", nil)
}

func (v *MonitorView) OnPacket(name string, event any) {
	if packet, ok := event.(sniffing.Packet); ok {
		fmt.Fprintf(v, "TCP FROM %s to %s\n", packet.Source, packet.Target)
	}
}
