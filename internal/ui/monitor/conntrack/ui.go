package conntrack

import (
	"fmt"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/conntrack"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/rivo/tview"
)

type ConntrackView struct {
	*tview.TextView
}

func NewContrackView() *ConntrackView {
	textView := tview.NewTextView()
	//textView.SetDynamicColors(true)
	textView.SetBorder(true).SetTitle("conntrack")
	v := ConntrackView{textView}

	state.Subscribe("monitor:conntrack", v.OnReceive)

	return &v
}

func (v *ConntrackView) OnReceive(name string, event any) {
	if item, ok := event.(conntrack.ConnectionTracked); ok {
		fmt.Fprintf(v,
			"%s:%d <=> %s:%d\n",
			item.Origin.SrcIP,
			item.Origin.SrcPort,
			item.Origin.DstIP,
			item.Origin.DstPort,
		)
	}
}
