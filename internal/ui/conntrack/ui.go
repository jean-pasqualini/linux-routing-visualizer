package conntrack

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/conntrack"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/k0kubun/pp"
	"github.com/rivo/tview"
)

type ConntrackView struct {
	*tview.TextView
}

func NewConntrackView() *ConntrackView {
	textView := tview.NewTextView().SetDynamicColors(true)
	v := &ConntrackView{
		textView,
	}

	state.Subscribe("conntrack:response", v.OnConntrackResponse)

	return v
}

func (v *ConntrackView) OnTabShow() {
	state.Dispatch("conntrack:request", nil)
}

func (v *ConntrackView) OnConntrackResponse(name string, event any) {
	v.Clear()
	if connList, ok := event.([]conntrack.ConnectionTracked); ok {
		pp.Fprint(tview.ANSIWriter(v), connList)
	}
}
