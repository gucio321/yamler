package widget

import (
	"fmt"
	"github.com/gucio321/yamler/pkg/workflow"
)

func GetOSs() []string {
	return []string{workflow.OSWindows, workflow.OSUbuntu, workflow.OSMacOS}
}

func (w *Widget) jobRunsOnID(jobName string) string {
	return fmt.Sprintf("JobRunsOn%v%s", w.id, jobName)
}

func ToStrSlice[T ~string](in []T) []string {
	out := make([]string, len(in))
	for i, v := range in {
		out[i] = string(v)
	}
	return out
}
