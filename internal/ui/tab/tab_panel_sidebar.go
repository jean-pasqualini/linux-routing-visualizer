package tab

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

type TabPanelSidebar struct {
	*tview.Box
	panel *TabPanel
	zones []int
}

func NewTabPanelSidebar(pages *tview.Pages) *TabPanel {
	v := NewTabPanel(pages, TabLayoutSidebar)
	v.Box.SetBackgroundColor(tcell.ColorDarkBlue)

	return v
}

func newSidebar(panel *TabPanel) *TabPanelSidebar {
	return &TabPanelSidebar{
		Box:   tview.NewBox(),
		panel: panel,
	}
}

func (v *TabPanelSidebar) Draw(screen tcell.Screen) {
	v.Box.DrawForSubclass(screen, v)

	if len(v.panel.tabNames) == 0 {
		return
	}

	v.drawTabBar(screen)
}

func (v *TabPanelSidebar) OnAddPage(name string) {
	y := 2
	if lenght := len(v.zones); lenght > 0 {
		y = v.zones[lenght-1] + 3
	}

	v.zones = append(v.zones, y)
}

func (v *TabPanelSidebar) drawTabBar(screen tcell.Screen) {
	x, y, _, _ := v.GetInnerRect()
	maxW := 30

	border := tcell.ColorMediumPurple
	activeBorder := tcell.ColorMediumVioletRed
	//text := tcell.ColorWhite
	//inactive := tcell.ColorGray

	curY := y + 1

	activeName := v.panel.getActiveName()
	for _, name := range v.panel.GetPageNames() {
		label := "      " + name + " "
		tabW := runewidth.StringWidth(label) + 2
		active := name == activeName

		bcol := border
		if active {
			bcol = activeBorder
		}

		tview.Print(screen, "─"+repeat("─", tabW-2)+"╮", x, curY, maxW, tview.AlignLeft, bcol)
		tview.Print(screen, " "+label+"│", x, curY+1, maxW, tview.AlignLeft, bcol)
		tview.Print(screen, "─"+repeat("─", tabW-2)+"╯", x, curY+2, maxW, tview.AlignLeft, bcol)

		curY += 3
	}
}

// MouseHandler returns the mouse handler for this primitive.
func (v *TabPanelSidebar) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		_, y := event.Position()

		if !v.InRect(event.Position()) {
			return false, nil
		}

		if action == tview.MouseLeftDown {
			_, innerY, _, _ := v.GetInnerRect()
			relativeToInner := y - innerY

			for index, zone := range v.zones {
				if relativeToInner >= zone-1 && relativeToInner <= zone+1 {
					pagesNames := v.panel.GetPageNames()
					v.panel.indexTab = uint8(index)
					v.panel.SwitchToPage(pagesNames[index])
				}
			}
		}

		return
	}
}
