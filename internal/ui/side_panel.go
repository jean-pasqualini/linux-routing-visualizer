package ui

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/diagram"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/tab"
	"github.com/rivo/tview"
)

func NewSidePanel() tview.Primitive {

	newNode := func(title string) *diagram.Node {
		return &diagram.Node{X: 15, W: 15, H: 5, Title: title}
	}

	buildCanvas := func(names []string) *diagram.DiagramCanvas {
		canvas := diagram.NewDiagramCanvas(80, 500)

		for _, name := range names {
			canvas.AddNode(newNode(name))
		}

		return canvas
	}

	/**
	tables := tview.NewTextView().
		SetScrollable(true).
		SetWrap(false). // important pour ANSI
		SetDynamicColors(true)
	*/
	pages := tview.NewPages()
	pages.AddPage("Tables", buildCanvas([]string{"Raw", "Mangle", "Nat", "Filter", "Security"}), true, true)
	pages.AddPage("Inbound", buildCanvas([]string{"PREROUTING", "INPUT"}), true, true)
	pages.AddPage("Forward", buildCanvas([]string{"PREROUTING", "FORWARD", "POSTROUTING"}), true, true)
	pages.AddPage("Outbound", buildCanvas([]string{"OUTPUT", "POSTROUTING"}), true, true)

	tabContainer := tab.NewTabPanelHozitonal(pages).SetOnSelect(func(name string) {
		pages.SwitchToPage(name)
	})

	return tabContainer
}
