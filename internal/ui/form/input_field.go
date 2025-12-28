package form

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func NewDesignInputField() *DesignInputField {
	input := tview.NewInputField()
	input.SetBorder(true)
	return &DesignInputField{
		InputField: input,
	}
}

type DesignInputField struct {
	*tview.InputField
}

func (i *DesignInputField) GetFieldHeight() int {
	return 3
}

// Draw draws this primitive onto the screen.
func (i *DesignInputField) Draw(screen tcell.Screen) {
	i.InputField.DrawForSubclass(screen, i)

	x, y, w, _ := i.GetInnerRect()

	//tview.Print(screen, strings.Repeat("=", w), x, y, w, tview.AlignLeft, tcell.ColorPurple)
	//tview.Print(screen, strings.Repeat("=", w), x, y+2, w, tview.AlignLeft, tcell.ColorPurple)

	i.InputField.SetRect(x, y, w, 1)
	//i.InputField.SetBackgroundColor(tcell.ColorLightGrey)
	i.InputField.Draw(screen)
}
