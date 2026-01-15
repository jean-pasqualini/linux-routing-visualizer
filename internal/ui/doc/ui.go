package doc

import (
	"strings"

	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/iptable"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/diagram"
	"github.com/rivo/tview"
)

func NewDocumentationView() tview.Primitive {

	newNode := func(title string) *diagram.Node {
		return &diagram.Node{X: 15, W: 15, H: 5, Title: title}
	}

	buildCanvas := func(names []string, title string) *diagram.DiagramCanvas {
		canvas := diagram.NewDiagramCanvas(80, 500)

		for _, name := range names {
			canvas.AddNode(newNode(name))
		}

		canvas.SetBorder(true)
		canvas.SetTitle(title)

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

	container := tview.NewFlex()
	container.AddItem(buildCanvas(listFromTables(iptable.TablesList[:]), "tables"), 0, 1, false)
	container.AddItem(buildCanvas(listFromChains(iptable.InboundChaining[:]), "inbound"), 0, 1, false)
	container.AddItem(buildCanvas(listFromChains(iptable.ForwardChaining[:]), "forward"), 0, 1, false)
	container.AddItem(buildCanvas(listFromChains(iptable.OutboundChaining[:]), "outbound"), 0, 1, false)

	return container
}
