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
		showUnmatched:      true,
	}

	state.Subscribe("simulator_result", rView.showResult)
	state.Subscribe("simulator:log", rView.addLog)

	return rView
}

type ResultView struct {
	*tab.TabPanelHorizontal
	treeView      *tview.TreeView
	logView       *tview.TextView
	showUnmatched bool
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

func (v *ResultView) chainIcon(chain simulator.SimulatorNetfilterChain) string {
	if chain.Decision == "DROP" {
		return "🕳️"
	}
	if chain.Decision == "REJECT" {
		return "🧱"
	}
	if chain.Decision == "ACCEPT" {
		return "✅"
	}
	if chain.Decision == "NONE" {
		return "➖"
	}
	if chain.Decision == "RETURN" {
		return "↩️"
	}

	return "🤷"
}

func (v *ResultView) ruleIcon(rule simulator.SimulatorResultRule) string {
	if !rule.Matched {
		return "🙈"
	}
	if rule.Action == "JUMP" {
		return "🦘"
	}
	if rule.Action == "TRACE" {
		return "🕳️"
	}
	if rule.Action == "DNAT" {
		return "🎯"
	}
	if rule.Action == "MASQUERADE" {
		return "🎭️"
	}
	if rule.Action == "ACCEPT" {
		return "✅"
	}
	if rule.Action == "DROP" {
		return "🕳️"
	}
	if rule.Action == "REJECT" {
		return "🧱"
	}
	if rule.Action == "RETURN" {
		return "↩️"
	}
	return "🤷"
}

func (v *ResultView) processNode(node *tview.TreeNode, events []simulator.SimulatorEvent) {
	for _, event := range events {
		if incoming, ok := event.(simulator.SimulatorIncomingInterface); ok {
			node.SetText("🖥️ packet in interface " + incoming.Interface)
		}
		if chain, ok := event.(simulator.SimulatorNetfilterChain); ok {
			chainIcon := v.chainIcon(chain)
			// 🔗
			chainNode := v.createNode("📍 " + chain.Name + " " + chainIcon)
			for _, rule := range chain.Rules {
				icon := v.ruleIcon(rule)
				color := "#39FF14::b"
				if !rule.Matched {
					color = "lightgrey:-:-"
				}
				outText := fmt.Sprintf("%s [%s]%s[white:-:-]", icon, color, rule.Raw)
				for _, mismatch := range rule.Mismatches {
					outText = text.TagColorRegions("#6D071A::b", color, outText, mismatch.Raw+" ")
				}
				ruleNode := v.createNode(outText)
				if rule.JumpChain != nil {
					v.processNode(ruleNode, []simulator.SimulatorEvent{*rule.JumpChain})
				}
				chainNode.AddChild(ruleNode)
			}
			chainNode.CollapseAll()
			node.AddChild(chainNode)
		}
		if route, ok := event.(simulator.SimulatorNetrouting); ok {
			node.AddChild(v.createNode("🧭 " + route.RouteType))
		}
	}
}

func (v *ResultView) SetResult(result simulator.SimulatorResult) {
	rootNode := tview.NewTreeNode("result")
	v.processNode(rootNode, result.Events)

	v.treeView.SetRoot(rootNode)
}
