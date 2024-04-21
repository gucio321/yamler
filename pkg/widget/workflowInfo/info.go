package workflowInfo

import "gopkg.in/yaml.v3"

type Info struct {
	Capture     bool
	Done        bool
	SearchError string

	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Inputs      map[string]struct {
		Description string `yaml:"description"`
		Default     string `yaml:"default"`
	} `yaml:"inputs"`
}

func Unmarshal(data []byte) *Info {
	info := &Info{}
	yaml.Unmarshal(data, info)
	return info
}
