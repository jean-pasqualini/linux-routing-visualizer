package simulator

import (
	"fmt"

	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/rivo/tview"
)

type FormEventTarget struct {
	Device string
	IP     string
	Port   string
}

type FormEvent struct {
	Source   FormEventTarget
	Target   FormEventTarget
	Protocol string
}

type SimulatorResultRuleEvent struct {
	Raw     string
	Matched bool
}
type SimulatorResultChainEvent struct {
	Name     string
	Decision string
}
type SimulatorResultEvent struct {
	Request FormEvent
	Rules   []SimulatorResultRuleEvent
	Chains  []SimulatorResultChainEvent
}

func NewSimulatorPanel() tview.Primitive {

	flex := tview.NewFlex()
	c := simulatorPanel{
		Flex:       flex,
		resultView: NewResultView(),
		logView:    tview.NewTextView(),
	}

	flex.SetDirection(tview.FlexRow)
	flex.AddItem(NewForm(), 18, 0, false)
	flex.AddItem(c.resultView, 0, 1, false)
	flex.AddItem(c.logView, 0, 1, false)

	state.Subscribe("simulator_result", c.showResult)
	state.Subscribe("logger", c.addLog)

	return c
}

type simulatorPanel struct {
	*tview.Flex
	resultView *ResultView
	logView    *tview.TextView
}

func (s *simulatorPanel) addLog(name string, event any) {
	w := s.logView.BatchWriter()
	defer w.Close()
	if event, ok := event.(string); ok {
		fmt.Fprintln(w, event)
	}
}

func (s *simulatorPanel) showResult(name string, event any) {
	if event, ok := event.(SimulatorResultEvent); ok {
		s.resultView.SetResult(event)
	}
}
