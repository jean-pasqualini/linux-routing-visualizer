package simulator

import (
	"fmt"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/simulator"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/text"
	"github.com/mattn/go-runewidth"
	"io"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/tab"
	"github.com/rivo/tview"
)

func NewResultView() *ResultView {
	textView := tview.NewTextView()
	textView.SetBackgroundColor(tcell.ColorDarkBlue)
	textView.SetDynamicColors(true)
	textView.SetBorder(true)
	textView.SetTextColor(tcell.ColorWhite)
	textView.SetBorderPadding(0, 0, 2, 2)
	logView := tview.NewTextView()
	pages := tab.NewTabPanelHozitonal(tview.NewPages())
	pages.AddPage("Rules", textView, true, true)
	pages.AddPage("Logs", logView, true, false)
	rView := &ResultView{
		TabPanelHorizontal: pages,
		textView:           textView,
		logView:            logView,
	}

	state.Subscribe("simulator_result", rView.showResult)
	state.Subscribe("logger", rView.addLog)

	return rView
}

type ResultView struct {
	*tab.TabPanelHorizontal
	textView *tview.TextView
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

func (v *ResultView) drawBox(w io.Writer, txt string) {
	size := runewidth.StringWidth(text.RemoveSquareBrackets(txt))
	fmt.Fprintf(w, "%s\n", "╭"+strings.Repeat("─", size+2)+"╮")
	fmt.Fprintf(w, "│ %s │\n", txt)
	fmt.Fprintf(w, "%s\n", "╰"+strings.Repeat("─", size+2)+"╯")
}

func (v *ResultView) SetResult(result simulator.SimulatorResult) {
	w := v.textView.BatchWriter()
	defer w.Close()
	w.Clear()
	for _, event := range result.Events {
		if incoming, ok := event.(simulator.SimulatorIncomingInterface); ok {
			v.drawBox(w, "Incoming packet in interface "+incoming.Interface)
		}
		if chain, ok := event.(simulator.SimulatorNetfilterChain); ok {
			v.drawBox(w, "Chain "+chain.Name+" "+chain.Decision)
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
				v.drawBox(w, outText)
			}
		}
		if route, ok := event.(simulator.SimulatorNetrouting); ok {
			v.drawBox(w, "Routing "+route.RouteType)
		}
	}

}
