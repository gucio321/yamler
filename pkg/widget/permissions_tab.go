package widget

import (
	"fmt"
	"github.com/AllenDang/giu"
	"github.com/gucio321/yamler/pkg/workflow"
)

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
