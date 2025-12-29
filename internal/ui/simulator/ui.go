package simulator

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/form/validator"
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
type SimulatorResultEvent struct {
	Request FormEvent
	Rules   []SimulatorResultRuleEvent
}

func NewSimulatorPanel() tview.Primitive {
	var formEvent FormEvent

	flex := tview.NewFlex()
	c := simulatorPanel{
		Flex:     flex,
		textView: tview.NewTextView(),
		logView:  tview.NewTextView(),
	}
	c.textView.SetBackgroundColor(tcell.ColorPowderBlue)
	c.textView.SetBorder(true)
	c.textView.SetBorderPadding(0, 0, 2, 2)

	formSource := tview.NewForm().
		AddInputField("IP", "", 20, validator.IpValidator, func(text string) {
			formEvent.Source.IP = text
		}).
		AddInputField("Port", "", 6, validator.PortValidator, func(text string) {
			formEvent.Source.Port = text
		}).
		AddInputField("Device", "", 6, validator.DeviceValidator, func(text string) {
			formEvent.Source.Device = text
		})
	formSource.SetFieldBackgroundColor(tcell.ColorWhite)
	formSource.SetBorder(true)
	formSource.SetTitle("Source")

	formTarget := tview.NewForm().
		AddInputField("IP", "", 20, validator.IpValidator, func(text string) {
			formEvent.Target.IP = text
		}).
		AddInputField("Port", "", 6, validator.PortValidator, func(text string) {
			formEvent.Target.Port = text
		}).
		AddInputField("Device", "", 6, validator.DeviceValidator, func(text string) {
			formEvent.Target.Device = text
		})
	formTarget.SetFieldBackgroundColor(tcell.ColorWhite)
	formTarget.SetBorder(true)
	formTarget.SetTitle("Target")

	formGeneral := tview.NewForm().
		AddDropDown("Protocol", []string{"TCP", "UDP"}, 0, func(option string, _ int) {
			formEvent.Protocol = option
		})
	formGeneral.SetFieldBackgroundColor(tcell.ColorWhite)
	formGeneral.SetBorder(true)
	formGeneral.SetTitle("General")

	simulateButton := tview.NewButton("SIMULATE").SetSelectedFunc(func() {
		state.Dispatch("logger", "je suis cliqué")
		state.Dispatch("simulator_request", formEvent)
	})

	formFlex := tview.NewFlex()
	formFlex.AddItem(formSource, 0, 1, false)
	formFlex.AddItem(formTarget, 0, 1, false)

	flex.SetDirection(tview.FlexRow)
	flex.AddItem(formGeneral, 5, 0, false)
	flex.AddItem(formFlex, 9, 0, false)
	flex.AddItem(simulateButton, 3, 0, false)
	flex.AddItem(c.textView, 0, 1, false)
	flex.AddItem(c.logView, 0, 1, false)

	state.Subscribe("simulator_result", c.showResult)
	state.Subscribe("logger", c.addLog)

	return c
}

type simulatorPanel struct {
	*tview.Flex
	textView *tview.TextView
	logView  *tview.TextView
}

func (s *simulatorPanel) addLog(name string, event any) {
	w := s.logView.BatchWriter()
	defer w.Close()
	if event, ok := event.(string); ok {
		fmt.Fprintln(w, event)
	}
}

func (s *simulatorPanel) showResult(name string, event any) {
	w := s.textView.BatchWriter()
	defer w.Close()
	w.Clear()
	if event, ok := event.(SimulatorResultEvent); ok {
		for _, rule := range event.Rules {
			icon := "❌"
			if rule.Matched {
				icon = "✅"
			}
			fmt.Fprintf(w, "%s %s\n", icon, rule.Raw)
		}
	} else {
		fmt.Println(w, "nooooon")
	}
}
