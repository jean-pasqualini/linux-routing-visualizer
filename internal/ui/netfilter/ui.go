package netfilter

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/netfilter/arptable"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/netfilter/conntrack"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/netfilter/iptable"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/tab"
	"github.com/rivo/tview"
)

type NetfilterView struct {
	*tab.TabPanel
}

func NewNetfilterView() *NetfilterView {

	iptableView := iptable.NewIpTableView()
	conntrackView := conntrack.NewConntrackView()
	arpTableView := arptable.NewArpTableView()
	ebTableView := tview.NewBox()

	pages := tab.NewTabPanelSidebar(tview.NewPages())
	pages.AddPage("conntrack", conntrackView, true, false)
	pages.AddPage("iptable", iptableView, true, false)
	pages.AddPage("arptable", arpTableView, true, false)
	pages.AddPage("ebtable", ebTableView, true, false)
	return &NetfilterView{
		TabPanel: pages,
	}
}
