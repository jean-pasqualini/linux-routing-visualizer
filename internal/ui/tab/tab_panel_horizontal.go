package tab

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/tview"
)

type TabPanelHorizontal struct {
	*tview.Box
	indexTab uint8
	tabNames []string
	zones    [][2]int
	pages    *tview.Pages
}

func NewTabPanelHozitonal(pages *tview.Pages) *TabPanelHorizontal {
	v := &TabPanelHorizontal{
		Box:      tview.NewBox(),
		pages:    pages,
		zones:    [][2]int{},
		indexTab: 0,
		tabNames: []string{},
	}

	v.SetBorder(false)
	v.SetBorderPadding(0, 0, 0, 0)
	return v
}

func (v *TabPanelHorizontal) AddPage(name string, item tview.Primitive, resize, visible bool) *tview.Pages {
	begin := 1
	if lenght := len(v.zones); lenght > 0 {
		begin = v.zones[lenght-1][1] + 1
	}

	end := begin + runewidth.StringWidth(name) + 3

	v.tabNames = append(v.tabNames, name)
	v.zones = append(v.zones, [2]int{begin, end})
	pages := v.pages.AddPage(name, item, resize, visible)
	pages.SendToBack(name)
	return pages
}

func (v *TabPanelHorizontal) Draw(screen tcell.Screen) {
	v.Box.DrawForSubclass(screen, v)

	if len(v.tabNames) == 0 {
		return
	}

	v.drawTabBar(screen)
	v.drawPages(screen)
}

func (v *TabPanelHorizontal) drawPages(screen tcell.Screen) {
	x, y, w, h := v.GetInnerRect()

	if h <= 0 || w <= 0 {
		return
	}

	// Clamp si pas assez de place
	topH := 3
	bottomH := h - topH

	if bottomH > 0 {
		v.pages.SetRect(x, y+topH, w, bottomH)
		v.pages.Draw(screen)
	}
}

func (v *TabPanelHorizontal) getActiveName() string {
	return v.tabNames[v.indexTab]
}

func (v *TabPanelHorizontal) drawTabBar(screen tcell.Screen) {
	x, y, w, _ := v.GetInnerRect()

	border := tcell.ColorMediumPurple
	activeBorder := tcell.ColorMediumVioletRed
	//text := tcell.ColorWhite
	//inactive := tcell.ColorGray

	curX := x + 1

	tview.Print(screen, repeat("─", w), x, y+2, 200, tview.AlignLeft, border)

	activeName := v.getActiveName()
	for _, name := range v.pages.GetPageNames(false) {
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

func (v *TabPanelHorizontal) Focus(delegate func(p tview.Primitive)) {
	delegate(v.pages)
	return
}

func (v *TabPanelHorizontal) HasFocus() bool {
	return v.pages.HasFocus()
}

// MouseHandler returns the mouse handler for this primitive.
func (v *TabPanelHorizontal) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		x, _ := event.Position()

		if !v.InRect(event.Position()) {
			return false, nil
		}

		if v.pages.InRect(event.Position()) {
			consumed, capture = v.pages.MouseHandler()(action, event, setFocus)
			return
		}

		if action == tview.MouseLeftDown {
			innerX, _, _, _ := v.GetInnerRect()
			relativeToInner := x - innerX

			for index, zone := range v.zones {
				if relativeToInner >= zone[0] && relativeToInner <= zone[1] {
					pagesNames := v.pages.GetPageNames(false)
					v.indexTab = uint8(index)
					v.pages.SwitchToPage(pagesNames[index])
				}
			}
		}

		return
	}
}

func (v *TabPanelHorizontal) nextTab() {
	v.indexTab++
	if uint8(len(v.tabNames)) < v.indexTab+1 {
		v.indexTab = 0
	}
	v.pages.SwitchToPage(v.getActiveName())
}

// InputHandler returns the handler for this primitive.
func (v *TabPanelHorizontal) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if v.HasFocus() {
			switch event.Key() {
			case tcell.KeyTab:
				v.nextTab()
				return
			}
		}

		if handler := v.pages.InputHandler(); handler != nil {
			handler(event, setFocus)
			return
		}
	}
}

func repeat(s string, n int) string {
	out := ""
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}
