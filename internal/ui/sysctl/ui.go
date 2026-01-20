package sysctl

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/sysctl"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/k0kubun/pp"
	"github.com/rivo/tview"
)

type SysctlView struct {
	*tview.TextView
	visible bool
}

func NewSysctlView() *SysctlView {
	v := &SysctlView{
		TextView: tview.NewTextView().SetDynamicColors(true),
		visible:  false,
	}

	state.Subscribe("sysctl:response", v.OnResponse)

	return v
}

func (v *SysctlView) OnTabShow() {
	state.Dispatch("sysctl:request", nil)
}

func (v *SysctlView) OnTabHide() {

}

func (v *SysctlView) OnResponse(name string, event any) {
	v.Clear()
	if config, ok := event.(sysctl.SysctlConfig); ok {
		pp.Fprintln(tview.ANSIWriter(v), config)
	}
}
