package ui

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/iptable"
	ui "github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/simulator"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/tab"
	"github.com/rivo/tview"
)

type MainPanel struct {
	*tab.TabPanelHorizontal
}

func NewMainPanel() *MainPanel {

	simulatorView := ui.NewSimulatorPanel()
	iptableView := iptable.NewIpTableView()

	pages := tab.NewTabPanelHozitonal(tview.NewPages())
	pages.AddPage("simulator", simulatorView, true, true)
	pages.AddPage("monitoring", tview.NewBox(), true, true)
	pages.AddPage("iptable", iptableView, true, true)
	pages.AddPage("routing", tview.NewBox(), true, false)
	pages.AddPage("sockets", tview.NewBox(), true, false)
	pages.AddPage("devices", tview.NewBox(), true, false)
	pages.AddPage("sysctl", tview.NewBox(), true, false)
	pages.AddPage("ipvs", tview.NewBox(), true, false)

	panel := &MainPanel{
		pages,
	}

	return panel
}
