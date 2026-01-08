package arp

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/arp"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/k0kubun/pp"
	"github.com/rivo/tview"
)

type ArpView struct {
	*tview.TextView
}

func NewArpView() *ArpView {
	textView := tview.NewTextView().SetDynamicColors(true)
	v := &ArpView{
		textView,
	}

	state.Subscribe("arp:response", v.OnArpResponse)

	return v
}

func (v *ArpView) OnTabShow() {
	state.Dispatch("arp:request", nil)
}

func (v *ArpView) OnArpResponse(name string, event any) {
	v.Clear()
	if arpList, ok := event.([]arp.ArpEntry); ok {
		pp.Fprint(tview.ANSIWriter(v), arpList)
	}
}
