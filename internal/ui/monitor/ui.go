package monitor

import (
	"fmt"

	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/sniffing"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/rivo/tview"
)

type MonitorView struct {
	*tview.Grid
	sniffView *tview.TextView
	started   bool
}

func NewMonitorView() *MonitorView {

	sniffView := tview.NewTextView()
	sniffView.SetBorder(true).SetTitle("Sniffing")

	grid := tview.NewGrid()

	rows := []int{0, 0, 0}
	cols := []int{0, 0, 0}

	grid.SetRows(rows...).SetColumns(cols...)

	empty := func() tview.Primitive { return tview.NewBox().SetBorder(true) }

	grid.AddItem(sniffView, 0, 0, 1, 1, 0, 0, false)

	grid.AddItem(empty(), 0, 1, 1, 1, 0, 0, false)

	grid.AddItem(empty(), 0, 2, 1, 1, 0, 0, false)

	grid.AddItem(empty(), 1, 0, 1, 1, 0, 0, false)

	grid.AddItem(empty(), 1, 1, 1, 1, 0, 0, false)

	grid.AddItem(empty(), 1, 2, 1, 1, 0, 0, false)

	grid.AddItem(empty(), 2, 0, 1, 1, 0, 0, false)

	grid.AddItem(empty(), 2, 1, 1, 1, 0, 0, false)

	grid.AddItem(empty(), 2, 2, 1, 1, 0, 0, false)

	mview := &MonitorView{
		Grid:      grid,
		sniffView: sniffView,
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

func (v *MonitorView) OnTabHide() {

}

func (v *MonitorView) OnPacket(name string, event any) {
	if packet, ok := event.(sniffing.Packet); ok {
		fmt.Fprintf(v.sniffView, "TCP FROM %s to %s\n", packet.Source, packet.Target)
	}
}
