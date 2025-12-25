package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func PrintMultiline(
	screen tcell.Screen,
	text string,
	x, y, width int,
	align int,
	color tcell.Color,
) {
	lines := tview.WordWrap(text, width)

	for i, line := range lines {
		tview.Print(screen, line, x, y+i, width, align, color)
	}
}
