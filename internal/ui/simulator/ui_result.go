package simulator

import (
	"fmt"

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
	pages := tview.NewPages()
	pages.AddPage("Logs", logView, true, false)
	pages.AddPage("Rules", textView, true, true)
	rView := &ResultView{
		TabPanelHorizontal: tab.NewTabPanelHozitonal(pages),
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
	if event, ok := event.(SimulatorResultEvent); ok {
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

func (v *ResultView) SetResult(result SimulatorResultEvent) {
	w := v.textView.BatchWriter()
	defer w.Close()
	w.Clear()
	fmt.Fprintf(w, "Chains (%d) :\n", len(result.Chains))
	for _, chain := range result.Chains {
		fmt.Fprintf(w, "%s: final decision is %s\n", chain.Name, chain.Decision)
	}
	fmt.Fprintln(w, "Rules :")
	for _, rule := range result.Rules {
		icon := "🙈"
		color := "lightgrey"
		if rule.Matched {
			icon = "✅"
			color = "#39FF14::b"
		}
		fmt.Fprintf(w, "%s [%s]%s[white:-:-]\n", icon, color, rule.Raw)
	}
}
