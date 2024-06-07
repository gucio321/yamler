package widget

import (
	"fmt"
	"github.com/AllenDang/giu"
	"github.com/gucio321/yamler/pkg/workflow"
	"os"
)

const Signature = `# Generated with yamler: https://github.com/gucio321/yamler by @gucio321`

type Widget struct {
	id    string
	w     *workflow.Workflow
	token string // GitHub token
}

func Workflow(w *workflow.Workflow) *Widget {
	return &Widget{
		w:  w,
		id: giu.GenAutoID("##WorkflowWidget"),
	}
}

func (w *Widget) Token(token string) *Widget {
	w.token = token
	return w
}

func (w *Widget) Build() {
	giu.Layout{
		w.requestsStatus(),
		giu.Separator(),
		w.workflowHeader(),
		giu.TabBar().TabItems(
			giu.TabItem("On (triggers)").Layout(w.triggersTab()),
			giu.TabItem("Permissions").Layout(w.permissionsTab()),
			giu.TabItem("Jobs").Layout(w.jobsTab()),
		),
	}.Build()
}

func (w *Widget) requestsStatus() giu.Layout {
	s := w.GetState()

	return giu.Layout{
		giu.CSSTag(func() string {
			if s.InvalidToken {
				return "error-detected"
			}
			return "main"
		}()).To(
			giu.Row(
				giu.Labelf("API Limits: %d of %d", s.APILimits.Remaining, s.APILimits.Limit),
				giu.SmallButton("Refresh").OnClick(func() {
					if err := w.updateRequestsLimit(); err != nil {
						fmt.Println("Error while getting API limits:", err)
					}
				}),
			),
		),
	}
}

func (w *Widget) workflowHeader() giu.Layout {
	s := w.GetState()

	return giu.Layout{
		giu.Row(
			giu.Label("Name:"),
			giu.InputText(&s.workflow.Name).Size(200),
			giu.Button("Generate Code").OnClick(func() {
				s.code, _ = s.workflow.Marshal()
				if s.signature {
					s.code = fmt.Sprintf("%s\n%s", Signature, s.code)
				}
				giu.OpenPopup("Code output")
			}),
			giu.Checkbox("Signature", &s.signature),
			giu.Tooltip(`Add a comment on top of generated code with a link to this generator

I do not enforce that but I'd be greatful!`),
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
		giu.Dummy(0, 20),
	}
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
