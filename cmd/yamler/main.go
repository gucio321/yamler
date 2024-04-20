package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/AllenDang/giu"
	"github.com/gucio321/yamler/pkg/widget"
	"github.com/gucio321/yamler/pkg/workflow"
)

func main() {
	// read whole stdin
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	// parse data
	w, err := workflow.Unmarshal(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(w.Marshal())
	wnd := giu.NewMasterWindow("Yamler", 640, 480, 0)
	wnd.Run(func() {
		giu.SingleWindow().Layout(
			widget.Workflow(w),
		)
	})
}
