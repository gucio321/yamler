package main

import (
	"fmt"

	"github.com/AllenDang/giu"
	"github.com/gucio321/yamler/pkg/widget"
	"github.com/gucio321/yamler/pkg/workflow"
)

func main() {
	w := &workflow.Workflow{}
	w.On.PageBuild.EnableEmpty = workflow.FieldOff
	w.Jobs = make(map[string]*workflow.Job)
	w.Jobs["tes"] = &workflow.Job{
		Steps: []*workflow.Step{
			{
				Name: "step1",
				Uses: "actions/checkout@v2",
			},
			{
				Name: "step2",
				Uses: "actions/setup-node@v2",
				With: map[string]string{
					"node-version": "14",
				},
			},
		},
	}
	fmt.Println(w.Marshal())
	wnd := giu.NewMasterWindow("Yamler", 640, 480, 0)
	wnd.Run(func() {
		giu.SingleWindow().Layout(
			widget.Workflow(w),
		)
	})
}
