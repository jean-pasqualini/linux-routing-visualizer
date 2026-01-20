package ebtable

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/ebtable"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/k0kubun/pp"
	"github.com/rivo/tview"
)

type EbtableView struct {
	*tview.TextView
}

func NewEbtableView() *EbtableView {
	textView := tview.NewTextView().SetDynamicColors(true)
	v := EbtableView{
		textView,
	}

	state.Subscribe("ebtable:response", v.OnResponse)

	return &v
}

func (v *EbtableView) OnTabShow() {
	state.Dispatch("ebtable:request", nil)
}

func (v *EbtableView) OnTabHide() {

}

func (v *EbtableView) OnResponse(name string, event any) {
	v.Clear()
	if table, ok := event.(ebtable.Ebtable); ok {
		pp.Fprintln(tview.ANSIWriter(v), table)
	}
}
