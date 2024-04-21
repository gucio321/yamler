package widget

import (
	"fmt"
	"github.com/AllenDang/giu"
	"github.com/gucio321/yamler/pkg/widget/workflowInfo"
	"github.com/gucio321/yamler/pkg/workflow"
	"net/http"
	"sort"
	"strings"
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
			giu.TreeNode("Steps").Layout(
				giu.Custom(func() {
					giu.Separator().Build()
					for i := 0; i < len(job.Steps); i++ {
						i := i
						w.jobStep(i, jobName, job.Steps[i]).Build()
						giu.Separator().Build()
					}
				}),
				giu.Button("Add step").OnClick(func() {
					job.Steps = append(job.Steps, &workflow.Step{})
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

func (w *Widget) jobStep(stepIdx int, jobID string, step *workflow.Step) giu.Widget {
	s := w.GetState()
	return giu.Layout{
		giu.Row(
			giu.Label("Name:"),
			giu.InputText(&step.Name).Size(100),
			giu.Label("ID:"),
			giu.InputText(&step.Id).Size(100),
		),
		giu.TreeNodef("Uses (External Action)##uses%v%v%v", w.id, jobID, stepIdx).Layout(
			giu.Row(
				giu.Label("Uses (Action ID):"),
				giu.InputText(&step.Uses).Hint("owner/repo@version").OnChange(func() {
					step.With = make(map[string]string)
					SearchActionInputs(step.Uses, s)
				}),
			),
			giu.Labelf("Name: %s", s.actionDetails.GetByID(step.Uses).Name),
			giu.Labelf("Description: %s", s.actionDetails.GetByID(step.Uses).Description),
			giu.Custom(func() {
				// here we print table with inputs
				info := s.actionDetails.GetByID(step.Uses)
				if !info.Done {
					return
				}
				rows := make([]*giu.TableRowWidget, len(info.Inputs))
				keys := make([]string, 0)
				for key := range info.Inputs {
					keys = append(keys, key)
				}

				sort.Strings(keys)

				for i, key := range keys {
					i := i
					rows[i] = giu.TableRow(
						giu.Layout{
							giu.Label(key),
							giu.Tooltip(info.Inputs[key].Description),
						},
						giu.InputText(s.actionsWith.GetByID(fmt.Sprintf("%s%s%d%s%s", key, step.Uses, stepIdx, jobID, w.id))).OnChange(func() {
							step.With[key] = *s.actionsWith.GetByID(fmt.Sprintf("%s%s%d%s%s", key, step.Uses, stepIdx, jobID, w.id))
							if step.With[key] == "" {
								delete(step.With, key)
							}
						}).Hint(info.Inputs[key].Default),
					)
				}

				if len(rows) == 0 {
					return
				}

				giu.Table().Rows(rows...).Build()
			}),
		),
		giu.TreeNodef("Script##script%v%v%v", w.id, jobID, stepIdx).Layout(
			giu.InputTextMultiline(&step.Run).Size(-1, 100),
		),
	}
}

// NOTE: this can't use w.GetState!
func SearchActionInputs(name string, s *State) {
	if s.actionDetails.GetByID(name).Capture {
		return
	}

	s.actionDetails.GetByID(name).Capture = true

	go func() {
		// try to extract action details from GitHub
		// and save it as value of s.actionDetails.GetByID(step.Uses)

		// [ preview ]
		// date we're looking for are stored in a action.yaml file
		// in action's repository on GitHub (we suppose github.com for now)
		// url will be of form https://raw.githubusercontent.com/owner/repo/version/action.yml
		// 1. send GET request to url
		url := fmt.Sprintf("https://raw.githubusercontent.com/%s/action.yml", strings.ReplaceAll(name, "@", "/"))
		request, err := http.NewRequest("GET", url, nil)
		if err != nil { // TODO: we can show this error somehow
			return
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			return
		}
		// 2. read response body
		output := make([]byte, 0)
		// read all content of response.Body
		// into output
		for {
			buffer := make([]byte, 1024)
			n, err := response.Body.Read(buffer)
			output = append(output, buffer[:n]...)
			if err != nil {
				break
			}
		}

		info := workflowInfo.Unmarshal(output)

		// 3. put it into s.actionDetails.GetByID(step.Uses)
		info.Capture = true
		info.Done = true
		// Reset this because most probably old options does not apply
		*s.actionDetails.GetByID(name) = *info
	}()
}
