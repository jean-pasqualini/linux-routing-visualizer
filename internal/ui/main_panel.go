package ui

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/tab"
	"github.com/k0kubun/pp"
	"github.com/rivo/tview"
)

type MainPanel struct {
	*tab.TabPanelHorizontal
	parsedView *tview.TextView
	rawView    *tview.TextView
}

func NewMainPanel() *MainPanel {
	parsedView := tview.NewTextView().
		SetScrollable(true).
		SetWrap(true).
		SetDynamicColors(true)

	rawView := tview.NewTextView().
		SetScrollable(true).
		SetWrap(true).
		SetDynamicColors(true)

	pages := tview.NewPages().
		AddPage("raw", rawView, true, true).
		AddPage("parsed", parsedView, true, true)

	return &MainPanel{
		tab.NewTabPanelHozitonal(pages),
		parsedView,
		rawView,
	}
}

func (p *MainPanel) ShowTables(app *tview.Application, tables any, raw string) {
	app.QueueUpdate(func() {
		p.parsedView.SetText(pp.Sprint(tables))
		p.rawView.SetText(raw)
	})
}
