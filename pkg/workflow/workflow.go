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

type Workflow struct {
	Name string `yaml:"name,omitempty"`
	On   struct {
		Push struct {
			EnableEmpty FieldSwitch `yaml:"-"`
			Branches    []string    `yaml:"branches",omitempty`
		} `yaml:"push,omitempty"`
		Fork  struct{} `yaml:"fork,omitempty"`
		Label struct {
			Types []OnLabelType `yaml:"types,omitempty"`
		} `yaml:"label,omitempty"`
		Issues struct {
			Types []OnIssueType `yaml:"types,omitempty"`
		} `yaml:"issues,omitempty"`
		PageBuild struct {
			EnableEmpty FieldSwitch `yaml:"-"`
		} `yaml:"page_build,omitempty"`
		PullRequest struct {
			EnableEmpty FieldSwitch `yaml:"-"`
			Types       []string    `yaml:"types,omitempty"`
		} `yaml:"pull_request,omitempty"`
	} `yaml:"on"`
}

func (w *Workflow) Marshal() (string, error) {
	result, err := yaml.Marshal(w)
	return string(result), err
}

type poweredMap map[string]interface{}
