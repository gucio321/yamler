package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/AllenDang/giu"
	"github.com/gucio321/yamler/pkg/assets"
	"github.com/gucio321/yamler/pkg/widget"
	"github.com/gucio321/yamler/pkg/workflow"
)

const tokenInfo = `GitHub token is recommended to increase API rate limit.
In order to generate token:
- Go to https://github.com
- Click on your profile icon in the top right corner
- Go to settings -> Developer settings -> Personal access tokens
- Click on "Generate new token"
- Select "repo" scope and click "Generate token"
- Copy and write down the token
`

func main() {
	n := flag.Bool("n", false, "create empty workflow instead of reading from stdin")
	token := flag.String("token", "", tokenInfo)
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

	wnd := giu.NewMasterWindow("Yamler", 640, 480, 0)
	if err := giu.ParseCSSStyleSheet(assets.Style); err != nil {
		log.Fatal(err)
	}

	wnd.Run(func() {
		giu.SingleWindow().Layout(
			widget.Workflow(w).Token(*token),
		)
	})
}
