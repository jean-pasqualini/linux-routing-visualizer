package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/tab"
	"github.com/k0kubun/pp"
	"github.com/rivo/tview"
)

type MainPanel struct {
	*tview.Box
	parsedView *tview.TextView
	rawView    *tview.TextView
	tabPanel   *tab.TabPanelHorizontal
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

	tabPanel := tab.NewTabPanelHozitonal(pages)

	return &MainPanel{
		tview.NewBox().SetBorder(true).SetTitle("Output"),
		parsedView,
		rawView,
		tabPanel,
	}
}

// MouseHandler returns the mouse handler for this primitive.
func (f *MainPanel) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return f.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		if !f.InRect(event.Position()) {
			return false, nil
		}

		consumed, capture = f.tabPanel.MouseHandler()(action, event, setFocus)
		if consumed {
			return
		}

		return
	})
}

func (m *MainPanel) Focus(delegate func(p tview.Primitive)) {
	m.tabPanel.Focus(delegate)
}

func (p *MainPanel) ShowTables(app *tview.Application, tables any) {
	app.QueueUpdate(func() {
		fmt.Fprintf(p.parsedView, pp.Sprint(tables))
	})
}
func (p *MainPanel) Draw(screen tcell.Screen) {
	p.Box.DrawForSubclass(screen, p)
	x, y, w, h := p.Box.GetInnerRect()

	p.tabPanel.SetRect(x, y, w, h)
	p.tabPanel.Draw(screen)
}
