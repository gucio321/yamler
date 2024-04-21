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
