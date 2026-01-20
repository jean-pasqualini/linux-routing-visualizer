package tab

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/event"
	"github.com/rivo/tview"
)

type TabPanel struct {
	*tview.Flex
	indexTab    uint8
	tabSelector TabSelectorView
	tabNames    []string
	pages       *tview.Pages
}

type TabLayout uint8

const (
	TabLayoutHorizontal TabLayout = iota
	TabLayoutSidebar
)

// [][2]int{}
func NewTabPanel(pages *tview.Pages, layout TabLayout) *TabPanel {
	v := &TabPanel{
		pages:    pages,
		indexTab: 0,
		tabNames: []string{},
	}

	flex := tview.NewFlex()
	if layout == TabLayoutSidebar {
		v.tabSelector = newSidebar(v)
		flex.SetDirection(tview.FlexColumn)
		flex.AddItem(v.tabSelector.(tview.Primitive), 30, 0, false)
	}
	if layout == TabLayoutHorizontal {
		v.tabSelector = newTop(v)
		flex.SetDirection(tview.FlexRow)
		flex.AddItem(v.tabSelector.(tview.Primitive), 3, 0, false)
	}

	flex.AddItem(pages, 0, 1, true)

	v.Flex = flex

	v.SetBorder(false)
	v.SetBorderPadding(0, 0, 0, 0)
	return v
}

func (v *TabPanel) AddPage(name string, item tview.Primitive, resize, visible bool) *TabPanel {
	if v.tabSelector != nil {
		v.tabSelector.OnAddPage(name)
	}
	v.tabNames = append(v.tabNames, name)
	pages := v.pages.AddPage(name, item, resize, visible)
	pages.SendToBack(name)
	if s, ok := item.(event.TabEventSubscribe); ok && visible {
		s.OnTabShow()
	}
	return v
}

func (v *TabPanel) getActiveName() string {
	return v.tabNames[v.indexTab]
}

func (v *TabPanel) Focus(delegate func(p tview.Primitive)) {
	delegate(v.pages)
	return
}

func (v *TabPanel) HasFocus() bool {
	return v.pages.HasFocus()
}

func (v *TabPanel) SwitchToPage(name string) {
	v.pages.SwitchToPage(name)
	primitive := v.pages.GetPage(name)
	if s, ok := primitive.(event.TabEventSubscribe); ok {
		s.OnTabShow()
	}
}

func (v *TabPanel) nextTab() {
	v.indexTab++
	if uint8(len(v.tabNames)) < v.indexTab+1 {
		v.indexTab = 0
	}
	v.SwitchToPage(v.getActiveName())
}

func (v *TabPanel) GetPageNames() []string {
	return v.pages.GetPageNames(false)
}

func repeat(s string, n int) string {
	out := ""
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}
