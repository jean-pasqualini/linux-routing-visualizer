package footer

import (
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/linux/network/ns"
	"github.com/jeanpasqualini/linux-routing-visualizer/internal/ui/state"
	"github.com/rivo/tview"
)

type FooterView struct {
	*tview.Form
	list []ns.NamespaceSpec
}

func NewFooterView() tview.Primitive {
	form := FooterView{
		Form: tview.NewForm(),
		list: []ns.NamespaceSpec{},
	}
	form.SetTitle("Options")
	form.SetBorder(true)

	state.Subscribe("namespace:response", form.OnResponse)
	state.Dispatch("namespace:request", nil)

	return &form
}

func (v *FooterView) OnResponse(name string, event any) {
	if nsList, ok := event.(ns.NamespaceList); ok {
		choices := []string{"Host"}
		for _, item := range nsList {
			v.list = append(v.list, item)
			choices = append(choices, item.Name)
		}
		v.AddDropDown("namespace", choices, 0, func(option string, optionIndex int) {
			if optionIndex == 0 {
				state.Dispatch("namespace:change", "/proc/self/ns/net")
			} else {
				state.Dispatch("namespace:change", v.list[optionIndex-1].Path)
			}
		})
	}
}
