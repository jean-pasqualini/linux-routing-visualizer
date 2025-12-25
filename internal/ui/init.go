package ui

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/iptable"
	"github.com/rivo/tview"
)

func start(app *tview.Application) {
	ipt := iptable.NewBackend()
	_, _ = ipt.ListChains("aeaze")
	go func() {
		//mainPanel.ShowTables(app, tables)
	}()

	app.QueueUpdate(func() {
		//fmt.Fprintf(p.parsedView, pp.Sprint(tables))
	})
}
