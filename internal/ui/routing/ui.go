package routing

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/routing"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/k0kubun/pp"
	"github.com/rivo/tview"
)

type RoutingView struct {
	*tview.TextView
}

func NewRoutingView() *RoutingView {
	textView := tview.NewTextView()
	textView.SetScrollable(true).
		SetWrap(true).
		SetDynamicColors(true)
	rView := &RoutingView{
		TextView: textView,
	}

	state.Subscribe("routing:response", rView.OnResponse)

	return rView
}

func (v *RoutingView) OnTabShow() {
	state.Dispatch("routing:request", nil)
}

func (v *RoutingView) OnTabHide() {

}

func (v *RoutingView) OnResponse(name string, event any) {
	v.Clear()
	if config, ok := event.(routing.RoutingConfig); ok {
		pp.Fprint(tview.ANSIWriter(v), config)
	}
}
