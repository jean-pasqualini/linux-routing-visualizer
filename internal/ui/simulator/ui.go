package simulator

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/form/validator"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/rivo/tview"
	"strings"
)

type FormEvent struct {
	SourceIP   string
	TargetIP   string
	TargetPort string
	Protocol   string
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
	var sourceIP string
	var targetIP string
	var targetPort string
	var protocol string = "TCP"

	flex := tview.NewFlex()
	c := simulatorPanel{
		Flex:     flex,
		textView: tview.NewTextView(),
		logView:  tview.NewTextView(),
	}
	c.textView.SetBackgroundColor(tcell.ColorPowderBlue)
	c.textView.SetBorder(true)
	c.textView.SetBorderPadding(2, 2, 2, 2)

	form := tview.NewForm().
		AddInputField("Source IP", "", 20, validator.IpValidator, func(text string) {
			sourceIP = text
		}).
		AddInputField("Target IP", "", 20, validator.IpValidator, func(text string) {
			targetIP = text
		}).
		AddInputField("Target Port", "", 6, validator.PortValidator, func(text string) {
			targetPort = text
		}).
		AddDropDown("Protocol", []string{"TCP", "UDP"}, 0, func(option string, _ int) {
			protocol = option
		}).
		AddButton("Simulate", func() {
			state.Dispatch("simulator_request", FormEvent{
				SourceIP:   sourceIP,
				TargetIP:   targetIP,
				TargetPort: targetPort,
				Protocol:   protocol,
			})
		})
	form.SetFieldBackgroundColor(tcell.ColorWhite)
	form.SetBorder(true)

	flex.SetDirection(tview.FlexRow)
	flex.AddItem(form, 15, 0, true)
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
		fmt.Fprintf(w, strings.Repeat("=", 15)+"\n")
		fmt.Fprintf(w, "SRC: %s\n", event.Request.SourceIP)
		fmt.Fprintf(w, "DEST: %s\n", event.Request.TargetIP)
		fmt.Fprintf(w, "DPORT: %s\n", event.Request.TargetPort)
		fmt.Fprintf(w, "PROTOCOL: %s\n", event.Request.Protocol)
		fmt.Fprintf(w, strings.Repeat("=", 15)+"\n")

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
