package iptable

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/iptable"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/tab"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/textview"
	"github.com/k0kubun/pp"
	"github.com/rivo/tview"
)

type IpTableView struct {
	*tab.TabPanel
	parsedView *textview.TextViewSearchable
	rawView    *tview.TextView
	visible    bool
}

type IpTableResponse struct {
	Parsed map[string]iptable.Table
	Raw    string
}

func NewIpTableView() *IpTableView {

	parsedView := textview.NewTextViewSearchable()

	rawView := tview.NewTextView().
		SetScrollable(true).
		SetWrap(true).
		SetDynamicColors(true)

	pages := tab.NewTabPanelTop(tview.NewPages())
	pages.AddPage("raw", rawView, true, true)
	pages.AddPage("parsed", parsedView, true, true)

	v := &IpTableView{
		parsedView: parsedView,
		rawView:    rawView,
		TabPanel:   pages,
		visible:    false,
	}

	state.Subscribe("iptables:response", v.render)

	return v
}

func (p *IpTableView) OnTabShow() {
	state.Dispatch("iptables:request", nil)
}

func (v *IpTableView) OnTabHide() {

}

func (p *IpTableView) render(name string, event any) {
	if event, ok := event.(IpTableResponse); ok {
		p.parsedView.SetText(pp.Sprint(event.Parsed))
		p.rawView.SetText(event.Raw)
	} else {
		p.parsedView.SetText("nooooo")
	}
}
