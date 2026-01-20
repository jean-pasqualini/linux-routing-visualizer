package monitor

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/monitor/conntrack"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/monitor/sniff"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/rivo/tview"
)

type MonitorView struct {
	*tview.Grid
	started bool
}

func NewMonitorView() *MonitorView {

	sniffView := sniff.NewSniffView()
	conntrackView := conntrack.NewContrackView()

	grid := tview.NewGrid()

	rows := []int{0, 0}
	cols := []int{0, 0}

	grid.SetRows(rows...).SetColumns(cols...)

	empty := func() tview.Primitive { return tview.NewBox().SetBorder(true) }

	grid.AddItem(sniffView, 0, 0, 1, 1, 0, 0, false)

	grid.AddItem(conntrackView, 0, 1, 1, 1, 0, 0, false)

	grid.AddItem(empty(), 1, 0, 1, 1, 0, 0, false)

	grid.AddItem(empty(), 1, 1, 1, 1, 0, 0, false)

	mview := &MonitorView{
		Grid: grid,
	}

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
	state.Dispatch("monitor:stop", nil)
}
