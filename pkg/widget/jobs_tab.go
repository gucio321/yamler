package widget

import (
	"fmt"
	"github.com/AllenDang/giu"
	"github.com/gucio321/yamler/pkg/workflow"
	"sort"
)

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
			giu.Child().Layout(
				giu.Labelf("Name: %s", jobName),
				giu.Row(
					giu.Label("Runs on:"),
					giu.Combo(
						fmt.Sprintf("##JobRunsOn%v%d", w.id, i),
						func() string {
							if job.RunsOn == "" {
								return "--"
							}
							return osList[*s.dropdowns.GetByID(w.jobRunsOnID(jobName))]
						}(),
						osList, s.dropdowns.GetByID(w.jobRunsOnID(jobName)),
					).OnChange(func() {
						job.RunsOn = workflow.OS(osList[*s.dropdowns.GetByID(w.jobRunsOnID(jobName))])
					}),
				),

				giu.TreeNodef("Needs##%s", jobName).Layout(
					w.needs(jobName),
				),
				// render steps
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
		),
		)
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
		giu.Dummy(0, 10),
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
			giu.CSSTag("delete-button").To(
				giu.Button("Delete").OnClick(func() {
					s.workflow.Jobs[jobID].Steps = append(s.workflow.Jobs[jobID].Steps[:stepIdx], s.workflow.Jobs[jobID].Steps[stepIdx+1:]...)
				}),
			),
		),
		giu.Style().SetDisabled(step.Run != "").To(
			giu.TreeNodef("Uses (External Action)##uses%v%v%v", w.id, jobID, stepIdx).Layout(
				giu.Row(
					giu.Label("Uses (Action ID):"),
					giu.Custom(func() {
						info := s.actionDetails.GetByID(step.Uses)
						i := giu.Row(
							giu.InputText(&step.Uses).Hint("owner/repo@version").OnChange(func() {
								s.branchesList = nil
							}),
							giu.Custom(func() {
								if s.branchesList != nil && len(s.branchesList) > 0 {
									giu.Combo(
										fmt.Sprintf("##branches%s%d", jobID, stepIdx),
										s.branchesList[s.currentBranch],
										s.branchesList,
										&s.currentBranch,
									).OnChange(func() {
										step.With = make(map[string]string)
										w.SearchActionInputs(step.Uses, s)
										s.APILimits.Dec()
									}).Build()
									return
								}

								giu.Button("Search available branches").OnClick(func() {
									w.SearchActionBranches(step.Uses, s)
									s.APILimits.Dec()
								}).Build()

								giu.Tooltip("This can't be automated because of GitHub API limitations.\nTODO maybe possibility to add token in future").Build()
							}),
						)

						if info.Done && info.SearchError != "" {
							giu.Layout{
								giu.CSSTag("error-detected").To(
									i,
									giu.Tooltip(info.SearchError),
								),
							}.Build()

							return
						}

						i.Build()
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

					giu.Table().Rows(rows...).Size(-1, 200).Build()
				}),
				giu.Style().SetDisabled(step.Uses == "").To(
					giu.CSSTag("delete-button").To(
						giu.Button("Clear").Size(-1, 0).OnClick(func() {
							step.Uses = ""
							step.With = make(map[string]string)
						}),
					),
				),
			),
		),
		giu.Style().SetDisabled(step.Uses != "").To(
			giu.TreeNodef("Script##script%v%v%v", w.id, jobID, stepIdx).Layout(
				giu.InputTextMultiline(&step.Run).Size(-1, 100),
				giu.Style().SetDisabled(step.Run == "").To(
					giu.CSSTag("delete-button").To(
						giu.Button("Clear").Size(-1, 0).OnClick(func() {
							step.Run = ""
						}),
					),
				),
			),
		),
	}
}
