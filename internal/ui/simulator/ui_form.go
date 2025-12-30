package simulator

import (
	"github.com/gdamore/tcell/v2"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/form/validator"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/rivo/tview"
)

type Form struct {
	*tview.Flex
}

func NewForm() *Form {
	return &Form{
		Flex: BuildForm(),
	}
}

func BuildForm() *tview.Flex {
	var formEvent FormEvent

	formSource := tview.NewForm().
		AddInputField("IP", "", 20, validator.IpValidator, func(text string) {
			formEvent.Source.IP = text
		}).
		AddInputField("Port", "", 6, validator.PortValidator, func(text string) {
			formEvent.Source.Port = text
		}).
		AddInputField("Device", "", 10, validator.DeviceValidator, func(text string) {
			formEvent.Source.Device = text
		})
	formSource.SetFieldBackgroundColor(tcell.ColorWhite)
	formSource.SetFieldTextColor(tcell.ColorBlack)
	formSource.SetBorder(true)
	formSource.SetTitle("Source")

	formTarget := tview.NewForm().
		AddInputField("IP", "", 20, validator.IpValidator, func(text string) {
			formEvent.Target.IP = text
		}).
		AddInputField("Port", "", 6, validator.PortValidator, func(text string) {
			formEvent.Target.Port = text
		}).
		AddInputField("Device", "", 10, validator.DeviceValidator, func(text string) {
			formEvent.Target.Device = text
		})
	formTarget.SetFieldBackgroundColor(tcell.ColorWhite)
	formTarget.SetFieldTextColor(tcell.ColorBlack)
	formTarget.SetBorder(true)
	formTarget.SetTitle("Target")

	formGeneral := tview.NewForm().SetHorizontal(true).
		AddDropDown("Protocol", []string{"TCP", "UDP", "ICMP"}, 0, func(option string, _ int) {
			formEvent.Protocol = option
		}).
		AddDropDown("State", []string{"NONE", "RELATED", "ESTABLISHED"}, 0, func(option string, _ int) {
			formEvent.Protocol = option
		})
	formGeneral.SetFieldBackgroundColor(tcell.ColorWhite)
	formGeneral.SetFieldTextColor(tcell.ColorBlack)
	formGeneral.SetBorder(true)
	formGeneral.SetTitle("General")

	simulateButton := tview.NewButton("SIMULATE").SetSelectedFunc(func() {
		state.Dispatch("logger", "je suis cliqué")
		state.Dispatch("simulator_request", formEvent)
	})
	styleButton := tcell.StyleDefault.Background(tcell.ColorMediumPurple).Foreground(tcell.ColorWhite)
	simulateButton.SetStyle(styleButton)
	simulateButton.SetBackgroundColorActivated(tcell.ColorMediumVioletRed)
	simulateButton.SetLabelColorActivated(tcell.ColorWhite)

	formFlex := tview.NewFlex()
	formFlex.AddItem(formSource, 0, 1, false)
	formFlex.AddItem(formTarget, 0, 1, false)

	mainFlex := tview.NewFlex()
	mainFlex.SetDirection(tview.FlexRow)
	mainFlex.AddItem(formGeneral, 5, 0, false)
	mainFlex.AddItem(formFlex, 9, 0, false)
	mainFlex.AddItem(simulateButton, 3, 0, false)

	return mainFlex
}
