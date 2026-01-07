package log

import (
	"fmt"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/rivo/tview"
)

type LogPanel struct {
	*tview.TextView
}

func NewLogPanel() *LogPanel {
	p := &LogPanel{
		TextView: tview.NewTextView(),
	}
	p.SetBorder(true).SetTitle("application logs")

	state.Subscribe("app:log", p.OnLog)

	return p
}

func (l *LogPanel) OnLog(name string, event any) {
	if msg, ok := event.(string); ok {
		fmt.Fprintln(l, msg)
	}
}
