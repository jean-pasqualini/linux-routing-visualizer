package ui

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/iptable"
	ui "github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/simulator"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/tab"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/textview"
	"github.com/k0kubun/pp"
	"github.com/rivo/tview"
)

type MainPanel struct {
	*tab.TabPanelHorizontal
	parsedView *textview.TextViewSearchable
	rawView    *tview.TextView
}

type IpTableResponse struct {
	Parsed map[string]iptable.Table
	Raw    string
}

func NewMainPanel() *MainPanel {
	parsedView := textview.NewTextViewSearchable()

	rawView := tview.NewTextView().
		SetScrollable(true).
		SetWrap(true).
		SetDynamicColors(true)

	simulatorView := ui.NewSimulatorPanel()

	pages := tab.NewTabPanelHozitonal(tview.NewPages())
	pages.AddPage("simulator", simulatorView, true, true)
	pages.AddPage("raw", rawView, true, true)
	pages.AddPage("parsed", parsedView, true, true)

	panel := &MainPanel{
		pages,
		parsedView,
		rawView,
	}

	state.Subscribe("iptables:response", panel.render)
	state.Dispatch("iptables:request", nil)

	return panel
}

func (p *MainPanel) render(name string, event any) {
	if event, ok := event.(IpTableResponse); ok {
		p.parsedView.SetText(pp.Sprint(event.Parsed))
		p.rawView.SetText(event.Raw)
	} else {
		p.parsedView.SetText("nooooo")
	}
}
