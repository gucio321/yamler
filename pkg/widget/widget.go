package widget

import (
	"fmt"
	"sort"

	"github.com/AllenDang/giu"
	"github.com/gucio321/yamler/pkg/workflow"
)

type Widget struct {
	id string
	w  *workflow.Workflow
}

func Workflow(w *workflow.Workflow) *Widget {
	return &Widget{
		w:  w,
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
			giu.TabItem("Jobs").Layout(w.jobsTab()),
		),
	}.Build()
}

func (w *Widget) triggersTab() giu.Widget {
	s := w.GetState()
	return giu.Layout{
		giu.Table().Rows(
			giu.TableRow(
				w.stateCheckbox("Push", "push_enabled", &s.workflow.On.Push.EnableEmpty),
				giu.Layout{
					giu.Child().Layout(
						giu.Label("Branches:"),
						w.dynamicList(&s.workflow.On.Push.Branches, &s.workflow.On.Push.EnableEmpty),
					).Size(0, 100),
					giu.Child().Layout(
						giu.Label("Tags:"),
						w.dynamicList(&s.workflow.On.Push.Tags, &s.workflow.On.Push.EnableEmpty),
					).Size(0, 100),
				},
			),
			giu.TableRow(
				w.stateCheckbox("Fork", "fork_enabled", &s.workflow.On.Fork.EnableEmpty),
			),
			giu.TableRow(
				w.stateCheckbox("Label", "label_enabled", &s.workflow.On.Label.EnableEmpty),
			),
		),
	}
}

func (w *Widget) dynamicList(list *[]string, forceEnable *workflow.FieldSwitch) giu.Widget {
	// s := w.GetState()

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

func (w *Widget) jobsTab() giu.Widget {
	osList := []string{workflow.OSWindows, workflow.OSUbuntu, workflow.OSMacOS}

	s := w.GetState()
	tabItems := make([]*giu.TabItemWidget, 0)
	names := make([]string, 0)
	for name := range s.workflow.Jobs {
		names = append(names, name)
	}

	sort.Strings(names)

	for i, jobName := range names {
		i := i
		job := s.workflow.Jobs[jobName]
		_ = job
		tabItems = append(tabItems, giu.TabItemf("%s##%d", jobName, i).Layout(
			giu.Labelf("Name: %s", jobName),
			giu.Row(
				giu.Label("Runs on:"),
				giu.Combo(
					fmt.Sprintf("##JobRunsOn%v%d", w.id, i),
					func() string {
						if job.RunsOn == "" {
							return "--"
						}
						return osList[*s.dropdowns.GetByID(fmt.Sprintf("JobRunsOn%v%d", w.id, i))]
					}(),
					osList, s.dropdowns.GetByID(fmt.Sprintf("JobRunsOn%v%d", w.id, i)),
				).OnChange(func() {
					job.RunsOn = workflow.OS(osList[*s.dropdowns.GetByID(fmt.Sprintf("JobRunsOn%v%d", w.id, i))])
				}),
			),
		))
	}

	return giu.Layout{
		giu.Row(
			giu.Label("Name: "),
			giu.InputText(&s.newJobName),
			giu.Button("Add job").OnClick(func() {
				s.workflow.Jobs[s.newJobName] = &workflow.Job{}
			}).Disabled(func() bool {
				_, ok := s.workflow.Jobs[s.newJobName]
				return ok || s.newJobName == ""
			}()),
		),
		giu.TabBar().TabItems(tabItems...),
	}
}
