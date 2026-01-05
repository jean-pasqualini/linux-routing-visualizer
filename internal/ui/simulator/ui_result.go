package simulator

import (
	"fmt"

	"github.com/jeanpasqualini/linux-routing-visualizer/internal/simulator"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/text"

	"github.com/gdamore/tcell/v2"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/tab"
	"github.com/rivo/tview"
)

func NewResultView() *ResultView {
	treeView := tview.NewTreeView()
	treeView.SetBackgroundColor(tcell.ColorDarkBlue)
	treeView.SetBorder(true)
	treeView.SetBorderPadding(0, 0, 2, 2)
	treeView.SetSelectedFunc(func(node *tview.TreeNode) {
		children := node.GetChildren()
		if len(children) > 0 {
			node.SetExpanded(!node.IsExpanded())
		}
	})
	logView := tview.NewTextView()
	pages := tab.NewTabPanelHozitonal(tview.NewPages())
	pages.AddPage("Rules", treeView, true, true)
	pages.AddPage("Logs", logView, true, false)
	rView := &ResultView{
		TabPanelHorizontal: pages,
		treeView:           treeView,
		logView:            logView,
	}

	state.Subscribe("simulator_result", rView.showResult)
	state.Subscribe("logger", rView.addLog)

	return rView
}

type ResultView struct {
	*tab.TabPanelHorizontal
	treeView *tview.TreeView
	logView  *tview.TextView
}

func (s *ResultView) showResult(name string, event any) {
	if event, ok := event.(simulator.SimulatorResult); ok {
		s.SetResult(event)
	}
}

func (s *ResultView) addLog(name string, event any) {
	w := s.logView.BatchWriter()
	defer w.Close()
	if event, ok := event.(string); ok {
		fmt.Fprintln(w, event)
	}
}

func (v *ResultView) createNode(txt string) *tview.TreeNode {
	stl := tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorWhite)
	return tview.NewTreeNode(txt).SetTextStyle(stl)
}

func (v *ResultView) SetResult(result simulator.SimulatorResult) {
	rootNode := tview.NewTreeNode("result")
	for _, event := range result.Events {
		if incoming, ok := event.(simulator.SimulatorIncomingInterface); ok {
			rootNode.SetText("🖥️ packet in interface " + incoming.Interface)
		}
		if chain, ok := event.(simulator.SimulatorNetfilterChain); ok {
			chainNode := v.createNode("🔗 " + chain.Name + " " + chain.Decision)
			for _, rule := range chain.Rules {
				icon := "✅"
				color := "#39FF14::b"
				if !rule.Matched {
					icon = "🙈"
					color = "lightgrey:-:-"
				}
				outText := fmt.Sprintf("%s [%s]%s[white:-:-]", icon, color, rule.Raw)
				for _, mismatch := range rule.Mismatches {
					outText = text.TagColorRegions("#6D071A::b", color, outText, mismatch.Raw+" ")
				}
				chainNode.AddChild(v.createNode(outText))
			}
			chainNode.CollapseAll()
			rootNode.AddChild(chainNode)
		}
		if route, ok := event.(simulator.SimulatorNetrouting); ok {
			rootNode.AddChild(v.createNode("🗺️ " + route.RouteType))
		}
	}

	v.treeView.SetRoot(rootNode)
}
