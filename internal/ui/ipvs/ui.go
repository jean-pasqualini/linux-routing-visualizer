package ipvs

import (
	"fmt"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/ipvs"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/k0kubun/pp"
	"github.com/rivo/tview"
)

type IPVSView struct {
	*tview.Flex
	rawText *tview.TextView
	goText  *tview.TextView
}

type IPVSEvent struct {
	Raw string
	Go  []ipvs.IPVSService
}

func NewIPVSView() *IPVSView {
	rawText := tview.NewTextView()
	rawText.SetDynamicColors(true).SetBorder(true).SetTitle("Raw version")
	goText := tview.NewTextView()
	goText.SetDynamicColors(true).SetBorder(true).SetTitle("Go version")
	view := &IPVSView{
		Flex:    tview.NewFlex(),
		rawText: rawText,
		goText:  goText,
	}
	view.AddItem(rawText, 0, 1, false)
	view.AddItem(goText, 0, 1, false)

	state.Subscribe("ipvs:response", view.onDisplay)

	return view
}

func (v *IPVSView) OnTabShow() {
	state.Dispatch("ipvs:request", nil)
}

func (v *IPVSView) OnTabHide() {

}

func (v *IPVSView) onDisplay(name string, event any) {
	v.rawText.Clear()
	v.goText.Clear()
	if event, ok := event.(IPVSEvent); ok {
		fmt.Fprintln(v.rawText, event.Raw)
		fmt.Fprintln(tview.ANSIWriter(v.goText), pp.Sprint(event.Go))
	}
}
