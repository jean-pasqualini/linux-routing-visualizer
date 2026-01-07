package ui

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/iptable"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/ipvs"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/monitor"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/routing"
	ui "github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/simulator"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/socket"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/tab"
	"github.com/rivo/tview"
)

type MainPanel struct {
	*tab.TabPanelHorizontal
}

func NewMainPanel() *MainPanel {

	simulatorView := ui.NewSimulatorPanel()
	iptableView := iptable.NewIpTableView()
	routingView := routing.NewRoutingView()
	IPVSView := ipvs.NewIPVSView()
	socketView := socket.NewSocketView()
	monitorView := monitor.NewMonitorView()

	pages := tab.NewTabPanelHozitonal(tview.NewPages())
	pages.AddPage("simulator", simulatorView, true, true)
	pages.AddPage("monitoring", monitorView, true, false)
	pages.AddPage("iptable", iptableView, true, false)
	pages.AddPage("routing", routingView, true, false)
	pages.AddPage("sockets", socketView, true, false)
	pages.AddPage("interfaces", tview.NewBox(), true, false)
	pages.AddPage("sysctl", tview.NewBox(), true, false)
	pages.AddPage("ipvs", IPVSView, true, false)

	panel := &MainPanel{
		pages,
	}

	return panel
}
