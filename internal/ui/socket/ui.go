package socket

import (
	"fmt"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/socket"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/k0kubun/pp"
	"github.com/rivo/tview"
)

type SocketView struct {
	*tview.Flex
	humanText *tview.TextView
	goText    *tview.TextView
}

func NewSocketView() *SocketView {
	humanText := tview.NewTextView()
	humanText.SetBorder(true).SetTitle("Human version")
	goText := tview.NewTextView()
	goText.SetDynamicColors(true).SetBorder(true).SetTitle("Go version")
	view := &SocketView{
		Flex:      tview.NewFlex(),
		humanText: humanText,
		goText:    goText,
	}
	view.AddItem(goText, 0, 1, false)
	view.AddItem(humanText, 0, 1, false)

	state.Subscribe("socket:response", view.onDisplay)
	state.Dispatch("socket:request", nil)

	return view
}

func (v *SocketView) onDisplay(name string, event any) {
	w := tview.ANSIWriter(v.goText)
	fmt.Fprintln(w, pp.Sprint(event))
	if sDescList, ok := event.([]socket.SocketDesc); ok {
		for _, sDesc := range sDescList {
			fmt.Fprintf(v.humanText, "%s(%d) is listening on port %d\n", sDesc.Comm, sDesc.Port, sDesc.PID)
		}
	}
}
