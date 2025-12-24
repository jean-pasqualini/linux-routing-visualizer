package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MainPanel struct {
	*tview.Box
}

func NewMainPanel() *MainPanel {
	return &MainPanel{
		tview.NewBox().SetBorder(true).SetTitle("Diagram"),
	}
}
func (p *MainPanel) Draw(screen tcell.Screen) {
	p.Box.DrawForSubclass(screen, p)

	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	screen.SetContent(5, 5, 'A', nil, style)
}
