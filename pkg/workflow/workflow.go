package workflow

import "gopkg.in/yaml.v3"

type OnLabelType string

const (
	OnLabelCreated OnLabelType = "created"
	OnLabelEdited  OnLabelType = "edited"
)

type OnIssueType string

const (
	OnIssueCreated OnIssueType = "created"
	OnIssueLabeled OnIssueType = "labeled"
)

type FieldSwitch interface{}

var (
	FieldOn  FieldSwitch = new(struct{ _ bool })
	FieldOff FieldSwitch = nil
)

func BoolToFieldSwitch(b bool) FieldSwitch {
	if b {
		return FieldOn
	}

	return FieldOff
}

type Permission string

const (
	PermRead  Permission = "read"
	PermWrite Permission = "write"
	PermNone  Permission = "none"
)

type OS string

const (
	OSWindows = "windows-latest"
	OSUbuntu  = "ubuntu-latest"
	OSMacOS   = "macos-latest"
)

type Workflow struct {
	Name    string `yaml:"name,omitempty"`
	RunName string `yaml:"run_name,omitempty"`
	On      struct {
		Push struct {
			EnableEmpty FieldSwitch `yaml:"-"`
			Branches    []string    `yaml:"branch,omitempty"`
			Tags        []string    `yaml:"tags,omitempty"`
		} `yaml:"push,omitempty"`
		Fork struct {
			EnableEmpty FieldSwitch `yaml:"-"`
		} `yaml:"fork,omitempty"`
		Label struct {
			EnableEmpty FieldSwitch   `yaml:"-"`
			Types       []OnLabelType `yaml:"types,omitempty"`
		} `yaml:"label,omitempty"`
		Issues struct {
			EnableEmpty FieldSwitch   `yaml:"-"`
			Types       []OnIssueType `yaml:"types,omitempty"`
		} `yaml:"issues,omitempty"`
		PageBuild struct {
			EnableEmpty FieldSwitch `yaml:"-"`
		} `yaml:"page_build,omitempty"`
		PullRequest struct {
			EnableEmpty FieldSwitch `yaml:"-"`

			Branches        []string `yaml:"branches,omitempty"`
			BranchesIgnored []string `yaml:"branches-ignored,omitempty"`
			Tags            []string `yaml:"tags,omitempty"`
			TagsIgnored     []string `yaml:"tags-ignored,omitempty"`
			Paths           []string `yaml:"paths,omitempty"`
			PathsIgnored    []string `yaml:"paths-ignored,omitempty"`

			Types []string `yaml:"types,omitempty"`
		} `yaml:"pull_request,omitempty"`
		PullRequestTarget struct {
			EnableEmpty FieldSwitch `yaml:"-"`

			Branches        []string `yaml:"branches,omitempty"`
			BranchesIgnored []string `yaml:"branches-ignored,omitempty"`
			Tags            []string `yaml:"tags,omitempty"`
			TagsIgnored     []string `yaml:"tags-ignored,omitempty"`
			Paths           []string `yaml:"paths,omitempty"`
			PathsIgnored    []string `yaml:"paths-ignored,omitempty"`

			Types []string `yaml:"types,omitempty"`
		} `yaml:"pull_request_target,omitempty"`
		WorkflowDispatch struct {
			EnableEmpty FieldSwitch         `yaml:"-"`
			Inputs      map[string]struct{} `yaml:"inputs,omitempty` // TODO
		} `yaml:"workflow_dispatch,omitempty"`
	} `yaml:"on"`
	Permissions struct {
		Actions            Permission `yaml:"actions,omitempty"`
		Checks             Permission `yaml:"checks,omitempty"`
		Contents           Permission `yaml:"contents,omitempty"`
		Deployments        Permission `yaml:"deployments,omitempty"`
		IDToken            Permission `yaml:"id-token,omitempty"`
		Issues             Permission `yaml:"issues,omitempty"`
		Discussions        Permission `yaml:"discussions,omitempty"`
		Packages           Permission `yaml:"packages,omitempty"`
		Pages              Permission `yaml:"pages,omitempty"`
		PullRequests       Permission `yaml:"pull-requests,omitempty"`
		RepositoryProjects Permission `yaml:"repository-projects,omitempty"`
		SecurityEvents     Permission `yaml:"security-events,omitempty"`
		Statuses           Permission `yaml:"statuses,omitempty"`
	} `yaml:"permissions,omitempty"`
	Jobs map[string]*Job `yaml:"jobs"`
}

type Job struct {
	RunsOn OS       `yaml:"runs-on"`
	Needs  []string `yaml:"needs,omitempty"`
	Steps  []*Step  `yaml:"steps"`
}

type Step struct {
	Name string            `yaml:"name,omitempty"`
	Id   string            `yaml:"id,omitempty"`
	Uses string            `yaml:"uses,omitempty"`
	With map[string]string `yaml:"with,omitempty"`
	Run  string            `yaml:"run,omitempty"`
}

func NewWorkflow() *Workflow {
	return &Workflow{
		Jobs: map[string]*Job{},
	}
}

func (w *Workflow) Marshal() (string, error) {
	result, err := yaml.Marshal(w)
	return string(result), err
}

func Unmarshal(data []byte) (*Workflow, error) {
	w := &Workflow{}
	err := yaml.Unmarshal(data, w)
	return w, err
}

type poweredMap map[string]interface{}
