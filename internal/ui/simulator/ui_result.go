package simulator

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewResultView() *ResultView {
	textView := tview.NewTextView()

	textView.SetBackgroundColor(tcell.ColorDarkBlue)
	textView.SetDynamicColors(true)
	textView.SetBorder(true)
	textView.SetTextColor(tcell.ColorWhite)
	textView.SetBorderPadding(0, 0, 2, 2)
	return &ResultView{
		TextView: textView,
	}
}

type ResultView struct {
	*tview.TextView
}

func (v *ResultView) SetResult(result SimulatorResultEvent) {
	w := v.TextView.BatchWriter()
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
