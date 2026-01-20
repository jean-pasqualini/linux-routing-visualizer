package arptable

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/arptable"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/k0kubun/pp"
	"github.com/rivo/tview"
)

type ArpTableView struct {
	*tview.TextView
	visible bool
}

func NewArpTableView() *ArpTableView {
	textView := tview.NewTextView().SetDynamicColors(true)
	v := &ArpTableView{
		textView, false,
	}

	state.Subscribe("arptable:response", v.OnResponse)

	return v
}

func (v *ArpTableView) OnTabShow() {
	state.Dispatch("arptable:request", nil)
}

func (v *ArpTableView) OnTabHide() {

}

func (v *ArpTableView) OnResponse(name string, event any) {
	v.Clear()
	if arpList, ok := event.([]arptable.ArpEntry); ok {
		pp.Fprint(tview.ANSIWriter(v), arpList)
	}
}
