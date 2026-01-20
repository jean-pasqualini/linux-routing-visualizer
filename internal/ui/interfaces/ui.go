package interfaces

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/netdevice"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/k0kubun/pp"
	"github.com/rivo/tview"
)

type InterfacesView struct {
	*tview.TextView
	visible bool
}

func NewInterfacesView() *InterfacesView {
	textView := tview.NewTextView().SetDynamicColors(true)
	v := &InterfacesView{textView, false}

	state.Subscribe("netdevice:response", v.OnResponse)
	state.Subscribe("namespace:changed", v.OnNamespaceChanged)
	return v
}

func (v *InterfacesView) OnNamespaceChanged(name string, event any) {
	if v.visible {
		state.Dispatch("netdevice:request", nil)
	}
}

func (v *InterfacesView) OnTabShow() {
	v.visible = true
	state.Dispatch("netdevice:request", nil)
}

func (v *InterfacesView) OnTabHide() {
	v.visible = false
}

func (v *InterfacesView) OnResponse(name string, event any) {
	v.Clear()
	if list, ok := event.([]netdevice.NetDeviceSpec); ok {
		pp.Fprint(tview.ANSIWriter(v), list)
	}
}
