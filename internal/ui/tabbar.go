package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TabBar struct {
	*tview.Box
	tabs     []string
	active   int
	onSelect func(index int, name string) // <-- callback custom
}

func NewTabBar(tabs []string) *TabBar {
	return &TabBar{
		Box:  tview.NewBox(),
		tabs: tabs,
	}
}

func (t *TabBar) SetOnSelect(fn func(index int, name string)) *TabBar {
	t.onSelect = fn
	return t
}

func (t *TabBar) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return func(ev *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch ev.Key() {
		case tcell.KeyTab:
			t.active = (t.active + 1) % len(t.tabs)
			if t.onSelect != nil {
				t.onSelect(t.active, t.tabs[t.active])
			}
		case tcell.KeyBacktab: // Shift+Tab
			t.active = (t.active - 1 + len(t.tabs)) % len(t.tabs)
			if t.onSelect != nil {
				t.onSelect(t.active, t.tabs[t.active])
			}
		}
	}
}

func (t *TabBar) Draw(screen tcell.Screen) {
	t.Box.DrawForSubclass(screen, t)

	x, y, w, _ := t.GetRect()

	for i, name := range t.tabs {
		label := " " + name + " "
		if i == t.active {
			label = "[ " + name + " ]"
		}
		tview.Print(screen, label, x+i*10, y, w, tview.AlignLeft, tcell.ColorWhite)
	}
}
