package simulator

import (
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
	State    string
	Protocol string
}

func NewSimulatorPanel() tview.Primitive {

	flex := tview.NewFlex()
	c := simulatorPanel{
		Flex:       flex,
		resultView: NewResultView(),
	}

	flex.SetDirection(tview.FlexRow)
	flex.AddItem(NewForm(), 18, 0, false)
	flex.AddItem(c.resultView, 0, 1, false)

	return c
}

type simulatorPanel struct {
	*tview.Flex
	resultView *ResultView
}
