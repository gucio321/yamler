package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/AllenDang/giu"
	"github.com/gucio321/yamler/pkg/widget"
	"github.com/gucio321/yamler/pkg/workflow"
)

func main() {
	n := flag.Bool("n", false, "create empty workflow instead of reading from stdin")
	flag.Parse()

	var w *workflow.Workflow = workflow.NewWorkflow()

	if !*n {
		// read whole stdin
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}

		// parse data
		w, err = workflow.Unmarshal(data)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(w.On.Push)

	wnd := giu.NewMasterWindow("Yamler", 640, 480, 0)
	wnd.Run(func() {
		giu.SingleWindow().Layout(
			widget.Workflow(w),
		)
	})
}
