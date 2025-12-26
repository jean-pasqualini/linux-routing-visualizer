package ui

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/iptable"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/diagram"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/tab"
	"github.com/rivo/tview"
	"strings"
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

	listFromChains := func(chains []iptable.ChainType) []string {
		output := []string{}

		for _, chain := range chains {
			output = append(output, string(chain))
		}

		return output
	}

	listFromTables := func(tables []iptable.TableType) []string {
		output := []string{}

		for _, table := range tables {
			sTable := string(table)
			output = append(output, strings.ToUpper(sTable[:1])+sTable[1:])
		}

		return output
	}

	/**
	tables := tview.NewTextView().
		SetScrollable(true).
		SetWrap(false). // important pour ANSI
		SetDynamicColors(true)
	*/
	pages := tview.NewPages()
	pages.AddPage("Tables", buildCanvas(listFromTables(iptable.TablesList[:])), true, true)
	pages.AddPage("Inbound", buildCanvas(listFromChains(iptable.InboundChaining[:])), true, true)
	pages.AddPage("Forward", buildCanvas(listFromChains(iptable.ForwardChaining[:])), true, true)
	pages.AddPage("Outbound", buildCanvas(listFromChains(iptable.OutboundChaining[:])), true, true)

	tabContainer := tab.NewTabPanelHozitonal(pages)

	return tabContainer
}
