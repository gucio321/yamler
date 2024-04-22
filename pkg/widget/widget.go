package widget

import (
	"fmt"
	"github.com/AllenDang/giu"
	"github.com/gucio321/yamler/pkg/widget/workflowInfo"
	"github.com/gucio321/yamler/pkg/workflow"
	"net/http"
	"os"
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
				giu.Button("Print to STDOUT").OnClick(func() {
					fmt.Fprintln(os.Stderr, "--Generating to stdout--")
					fmt.Fprintln(os.Stdout, s.code)
				}).Size(-1, 0),
				giu.Button("Close").OnClick(func() {
					giu.CloseCurrentPopup()
				}).Size(-1, 0),
			),
		),
		giu.Row(
			giu.Label("Run Name (?):"),
			giu.Tooltip("").Layout(
				giu.Label(`The name for workflow runs generated from the workflow.
GitHub displays the workflow run name in the list of workflow runs
on your repository's "Actions" tab. If run-name is omitted or is only
whitespace, then the run name is set to event-specific information
for the workflow run. For example, for a workflow triggered by a push
or pull_request event, it is set as the commit message or the title of the
pull request.

This value can include expressions and can reference the github and inputs contexts.`),
			),

			giu.InputText(&s.workflow.RunName).Hint("e.g. Deploy to ${{ inputs.deploy_target }} by @${{ github.actor }}"),
		),
		giu.TabBar().TabItems(
			giu.TabItem("On (triggers)").Layout(w.triggersTab()),
			giu.TabItem("Permissions").Layout(w.permissionsTab()),
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

func (w *Widget) permissionsTab() giu.Widget {
	s := w.GetState()
	permissions := ToStrSlice([]workflow.Permission{workflow.PermNone, workflow.PermRead, workflow.PermWrite})
	rowsPresets := []struct {
		superMapID string
		field      *workflow.Permission
	}{
		{"actions", &s.workflow.Permissions.Actions},
		{"checks", &s.workflow.Permissions.Checks},
		{"contents", &s.workflow.Permissions.Contents},
		{"deployments", &s.workflow.Permissions.Deployments},
		{"idToken", &s.workflow.Permissions.IDToken},
		{"issues", &s.workflow.Permissions.Issues},
		{"discussions", &s.workflow.Permissions.Discussions},
		{"packages", &s.workflow.Permissions.Packages},
		{"pages", &s.workflow.Permissions.Pages},
		{"pullRequests", &s.workflow.Permissions.PullRequests},
		{"repositoryProjects", &s.workflow.Permissions.RepositoryProjects},
		{"securityEvents", &s.workflow.Permissions.SecurityEvents},
		{"statuses", &s.workflow.Permissions.Statuses},
	}

	return giu.Layout{
		giu.Label("If you specify the access for any of these scopes, all of those that are not specified are set to none."),
		giu.Table().Rows(func() []*giu.TableRowWidget {
			result := make([]*giu.TableRowWidget, 0)
			for _, row := range rowsPresets {
				row := row
				yield := giu.TableRow(
					giu.Label(row.superMapID),
					giu.Row(
						giu.Combo(
							fmt.Sprintf("##%s", row.superMapID),
							string(*row.field),
							permissions,
							s.dropdowns.GetByID(row.superMapID),
						).OnChange(func() {
							*row.field = workflow.Permission(permissions[*s.dropdowns.GetByID(row.superMapID)])
						}),
						giu.CSSTag("delete-button").To(
							giu.Button("Reset").OnClick(func() {
								*row.field = ""
							}),
						),
					),
				)
				result = append(result, yield)
			}
			return result
		}()...),
	}
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
						return osList[*s.dropdowns.GetByID(w.jobRunsOnID(jobName))]
					}(),
					osList, s.dropdowns.GetByID(w.jobRunsOnID(jobName)),
				).OnChange(func() {
					job.RunsOn = workflow.OS(osList[*s.dropdowns.GetByID(w.jobRunsOnID(jobName))])
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
						i := giu.InputText(&step.Uses).Hint("owner/repo@version").OnChange(func() {
							step.With = make(map[string]string)
							SearchActionInputs(step.Uses, s)
						})

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
			),
		),
		giu.Style().SetDisabled(step.Uses != "").To(
			giu.TreeNodef("Script##script%v%v%v", w.id, jobID, stepIdx).Layout(
				giu.InputTextMultiline(&step.Run).Size(-1, 100),
			),
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
		if err != nil {
			info := &workflowInfo.Info{
				Capture:     true,
				Done:        true,
				SearchError: err.Error(),
			}

			*s.actionDetails.GetByID(name) = *info
			return
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			info := &workflowInfo.Info{
				Capture:     true,
				Done:        true,
				SearchError: err.Error(),
			}

			*s.actionDetails.GetByID(name) = *info
			return
		}

		if response.StatusCode != 200 {
			info := &workflowInfo.Info{
				Capture: true,
				Done:    true,
			}

			switch response.StatusCode {
			case 404:
				info.SearchError = "Action not found"
			case 400:
				info.SearchError = "Invalid action name format"
			default:
				info.SearchError = fmt.Sprintf("Unexpected status code: %d", response.StatusCode)
			}

			*s.actionDetails.GetByID(name) = *info
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
