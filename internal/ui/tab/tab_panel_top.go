package tab

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

type TabPanelTop struct {
	*tview.Box
	panel *TabPanel
	zones [][2]int
}

func NewTabPanelTop(pages *tview.Pages) *TabPanel {
	v := NewTabPanel(pages, TabLayoutHorizontal)

	return v
}

func newTop(panel *TabPanel) *TabPanelTop {
	return &TabPanelTop{
		Box:   tview.NewBox(),
		zones: [][2]int{},
		panel: panel,
	}
}

func (v *TabPanelTop) OnAddPage(name string) {
	begin := 1
	if lenght := len(v.zones); lenght > 0 {
		begin = v.zones[lenght-1][1] + 1
	}

	end := begin + runewidth.StringWidth(name) + 3

	v.zones = append(v.zones, [2]int{begin, end})
}

func (v *TabPanelTop) Draw(screen tcell.Screen) {
	v.Box.DrawForSubclass(screen, v)

	if len(v.panel.tabNames) == 0 {
		return
	}

	v.drawTabBar(screen)
}

// MouseHandler returns the mouse handler for this primitive.
func (v *TabPanelTop) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		x, _ := event.Position()

		if !v.InRect(event.Position()) {
			return false, nil
		}

		if action == tview.MouseLeftDown {
			innerX, _, _, _ := v.GetInnerRect()
			relativeToInner := x - innerX

			for index, zone := range v.zones {
				if relativeToInner >= zone[0] && relativeToInner <= zone[1] {
					pagesNames := v.panel.GetPageNames()
					v.panel.indexTab = uint8(index)
					v.panel.SwitchToPage(pagesNames[index])
				}
			}
		}

		return
	}
}

func (v *TabPanelTop) drawTabBar(screen tcell.Screen) {
	x, y, w, _ := v.GetInnerRect()

	border := tcell.ColorMediumPurple
	activeBorder := tcell.ColorMediumVioletRed
	//text := tcell.ColorWhite
	//inactive := tcell.ColorGraysss

	curX := x + 1

	tview.Print(screen, repeat("─", w), x, y+2, w, tview.AlignLeft, border)

	activeName := v.panel.getActiveName()
	for _, name := range v.panel.GetPageNames() {
		label := " " + name + " "
		tabW := runewidth.StringWidth(label) + 2
		active := name == activeName

		bcol := border
		if active {
			bcol = activeBorder
		}

		tview.Print(screen, "╭"+repeat("─", tabW-2)+"╮", curX, y, w, tview.AlignLeft, bcol)
		tview.Print(screen, "│"+label+"│", curX, y+1, w, tview.AlignLeft, bcol)

		if active {
			tview.Print(screen, "┘"+repeat(" ", tabW-2)+"└", curX, y+2, w, tview.AlignLeft, bcol)
		} else {
			tview.Print(screen, "┴"+repeat("─", tabW-2)+"┴", curX, y+2, w, tview.AlignLeft, bcol)
		}

		curX += tabW
	}
}
