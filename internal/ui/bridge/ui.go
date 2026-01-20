package bridge

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/bridge"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/k0kubun/pp"
	"github.com/rivo/tview"
)

type BridgeView struct {
	*tview.TextView
	visible bool
}

func NewBridgeView() *BridgeView {
	textView := tview.NewTextView().SetDynamicColors(true)
	v := BridgeView{
		textView,
		false,
	}

	state.Subscribe("bridge:response", v.OnResponse)

	return &v
}

func (v *BridgeView) OnTabShow() {
	state.Dispatch("bridge:request", nil)
}

func (v *BridgeView) OnTabHide() {

}

func (v *BridgeView) OnResponse(name string, event any) {
	v.Clear()
	if list, ok := event.([]bridge.BridgeSpec); ok {
		pp.Fprintln(tview.ANSIWriter(v), list)
	}
}
