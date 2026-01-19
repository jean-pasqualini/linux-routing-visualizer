package ui

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/arp"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/bridge"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/doc"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/interfaces"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/ipvs"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/monitor"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/netfilter"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/routing"
	ui "github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/simulator"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/socket"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/sysctl"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/tab"
	"github.com/rivo/tview"
)

type MainPanel struct {
	*tab.TabPanel
}

func NewMainPanel() *MainPanel {

	simulatorView := ui.NewSimulatorPanel()
	routingView := routing.NewRoutingView()
	IPVSView := ipvs.NewIPVSView()
	socketView := socket.NewSocketView()
	monitorView := monitor.NewMonitorView()
	sysctlView := sysctl.NewSysctlView()
	arpView := arp.NewArpView()
	interfacesView := interfaces.NewInterfacesView()
	documentationView := doc.NewDocumentationView()
	netFilterView := netfilter.NewNetfilterView()
	bridgeView := bridge.NewBridgeView()

	pages := tab.NewTabPanelTop(tview.NewPages())
	pages.AddPage("simulator", simulatorView, true, true)
	pages.AddPage("monitoring", monitorView, true, false)
	pages.AddPage("netfilter", netFilterView, true, false)
	pages.AddPage("arp map", arpView, true, false)
	pages.AddPage("routing", routingView, true, false)
	pages.AddPage("sockets", socketView, true, false)
	pages.AddPage("interfaces", interfacesView, true, false)
	pages.AddPage("bridges", bridgeView, true, false)
	pages.AddPage("sysctl", sysctlView, true, false)
	pages.AddPage("ipvs", IPVSView, true, false)
	pages.AddPage("documentation", documentationView, true, false)

	panel := &MainPanel{
		pages,
	}

	return panel
}
