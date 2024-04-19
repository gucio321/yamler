package widget

import (
	"github.com/AllenDang/giu"
	"github.com/gucio321/yamler/pkg/workflow"
)

type Widget struct {
	id string
}

func Workflow() *Widget {
	return &Widget{
		id: giu.GenAutoID("##WorkflowWidget"),
	}
}

func (w *Widget) Build() {
	s := w.GetState()
	giu.Layout{
		giu.Row(
			giu.Label("Name:"),
			giu.InputText(&s.workflow.Name),
			giu.Button("Generate Code").OnClick(func() {
				s.code, _ = s.workflow.Marshal()
				giu.OpenPopup("Code output")
			}),
			giu.PopupModal("Code output").Layout(
				giu.Child().Layout(
					giu.InputTextMultiline(&s.code).Size(-1, -1),
				).Size(300, 300),
				giu.Button("Close").OnClick(func() {
					giu.CloseCurrentPopup()
				}).Size(-1, 0),
			),
		),
		giu.TabBar().TabItems(
			giu.TabItem("On (triggers)").Layout(w.triggersTab()),
		),
	}.Build()
}

func (w *Widget) triggersTab() giu.Widget {
	s := w.GetState()
	return giu.Layout{
		giu.Table().Rows(
			giu.TableRow(
				giu.Label("Push"),
				giu.Layout{
					giu.Checkbox("Enabled", s.toggles.GetByID("push_enabled")).OnChange(func() {
						s.workflow.On.Push.EnableEmpty = workflow.BoolToFieldSwitch(*s.toggles.GetByID("push_enabled"))
					}),
					giu.Child().Layout(
						giu.Label("Branches:"),
						giu.Custom(func() {
							for i := 0; i < len(s.workflow.On.Push.Branches); i++ {
								branch := s.workflow.On.Push.Branches[i]
								if branch == "" {
									s.workflow.On.Push.Branches = append(s.workflow.On.Push.Branches[:i], s.workflow.On.Push.Branches[i+1:]...)
									continue
								}

								giu.InputText(&s.workflow.On.Push.Branches[i]).Build()
							}

							tmp := ""
							giu.InputText(&tmp).Hint("Add branch").OnChange(func() {
								s.workflow.On.Push.Branches = append(s.workflow.On.Push.Branches, tmp)
								tmp = ""
							}).Build()
						}),
					),
				},
			),
		),
	}
}
