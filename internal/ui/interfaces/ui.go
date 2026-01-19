package interfaces

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/netdevice"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/k0kubun/pp"
	"github.com/rivo/tview"
)

type InterfacesView struct {
	*tview.TextView
}

func NewInterfacesView() *InterfacesView {
	textView := tview.NewTextView().SetDynamicColors(true)
	v := &InterfacesView{textView}

	state.Subscribe("netdevice:response", v.OnResponse)
	return v
}

func (v *InterfacesView) OnTabShow() {
	state.Dispatch("netdevice:request", nil)
}

func (v *InterfacesView) OnResponse(name string, event any) {
	v.Clear()
	if list, ok := event.([]netdevice.NetDeviceSpec); ok {
		pp.Fprint(tview.ANSIWriter(v), list)
	}
}
