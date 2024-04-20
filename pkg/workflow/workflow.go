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

type OS string

const (
	OSWindows = "windows-latest"
	OSUbuntu  = "ubuntu-latest"
	OSMacOS   = "macos-latest"
)

type Workflow struct {
	Name string `yaml:"name,omitempty"`
	On   struct {
		Push struct {
			EnableEmpty FieldSwitch `yaml:"-"`
			Branches    []string    `yaml:"branches,omitempty"`
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
			Types       []string    `yaml:"types,omitempty"`
		} `yaml:"pull_request,omitempty"`
	} `yaml:"on"`
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
