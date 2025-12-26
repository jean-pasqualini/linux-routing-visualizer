package ui

import "github.com/rivo/tview"

func NewSimulatorPanel() tview.Primitive {
	var sourceIP string
	var targetIP string
	var targetPort string
	var protocol string = "TCP"

	form := tview.NewForm().
		AddInputField("Source IP", "", 20, nil, func(text string) {
			sourceIP = text
		}).
		AddInputField("Target IP", "", 20, nil, func(text string) {
			targetIP = text
		}).
		AddInputField("Target Port", "", 6, nil, func(text string) {
			targetPort = text
		}).
		AddDropDown("Protocol", []string{"TCP", "UDP"}, 0, func(option string, _ int) {
			protocol = option
		}).
		AddButton("Simulate", func() {
			showResult(sourceIP, targetIP, targetPort, protocol)
		})

	return form
}

func showResult(sourceIP string, targetIP string, targetPort string, protocol string) {

}
