package conntrack

import (
	"fmt"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/conntrack"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/k0kubun/pp"
	"github.com/rivo/tview"
)

type ConntrackView struct {
	*tview.Flex
	goText    *tview.TextView
	humanText *tview.TextView
}

func NewConntrackView() *ConntrackView {
	flexView := tview.NewFlex()
	goTextView := tview.NewTextView().SetDynamicColors(true)
	humanTextView := tview.NewTextView().SetDynamicColors(true)
	v := &ConntrackView{
		flexView,
		goTextView,
		humanTextView,
	}

	v.AddItem(goTextView, 0, 1, false)
	v.AddItem(humanTextView, 0, 1, false)

	state.Subscribe("conntrack:response", v.OnConntrackResponse)

	return v
}

func (v *ConntrackView) OnTabShow() {
	state.Dispatch("conntrack:request", nil)
}

func (v *ConntrackView) OnConntrackResponse(name string, event any) {
	v.goText.Clear()
	v.humanText.Clear()
	if connList, ok := event.([]conntrack.ConnectionTracked); ok {
		pp.Fprint(tview.ANSIWriter(v.goText), connList)
		for _, conn := range connList {
			fmt.Fprintf(v.humanText, "%s -> %s\n", conn.Origin.SrcIP, conn.Origin.DstIP)
			fmt.Fprintf(v.humanText, "%s -> %s\n", conn.Return.SrcIP, conn.Return.DstIP)
			fmt.Fprintln(v.humanText, "---------")
		}
	}
}
