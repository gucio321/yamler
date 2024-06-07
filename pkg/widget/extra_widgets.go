package widget

import (
	"github.com/AllenDang/giu"
	"github.com/gucio321/yamler/pkg/workflow"
)

func (w *Widget) dynamicList(list *[]string, forceEnable *workflow.FieldSwitch) giu.Widget {
	return giu.Custom(func() {
		for i := 0; i < len(*list); i++ {
			branch := (*list)[i]
			if branch == "" {
				*list = append((*list)[:i], (*list)[i+1:]...)
				/*
					if len(*list) == 0 {
						*forceEnable = workflow.BoolToFieldSwitch(false)
						*s.toggles.GetByID("push_enabled") = false
					}
				*/
				continue
			}

			giu.InputText(&(*list)[i]).Build()
		}

		tmp := ""
		giu.InputText(&tmp).Hint("Add...").OnChange(func() {
			//*s.toggles.GetByID("push_enabled") = true
			//*forceEnable = workflow.BoolToFieldSwitch(true)
			*list = append((*list), tmp)
			tmp = ""
		}).Build()
	})
}

func (w *Widget) stateCheckbox(label, tag string, sw *workflow.FieldSwitch) giu.Widget {
	s := w.GetState()
	return giu.Checkbox(label, s.toggles.GetByID(tag)).OnChange(func() {
		*sw = workflow.BoolToFieldSwitch(*s.toggles.GetByID(tag))
	})
}
