package widget

import (
	"encoding/json"
	"fmt"
	"github.com/gucio321/yamler/pkg/widget/workflowInfo"
	"net/http"
	"strings"
)

// NOTE: this can't use w.GetState!
func (w *Widget) SearchActionInputs(name string, s *State) {
	if s.actionDetails.GetByID(name).Capture {
		return
	}

	s.actionDetails.GetByID(name).Capture = true

	go func() {
		// try to extract action details from GitHub
		// and save it as value of s.actionDetails.GetByID(step.Uses)

		// [ preview ]
		// date we're looking for are stored in a action.yaml file
		// in action's repository on GitHub (we suppose github.com for now)
		// url will be of form https://raw.githubusercontent.com/owner/repo/version/action.yml
		// 1. send GET request to url
		url := fmt.Sprintf("https://raw.githubusercontent.com/%s/action.yml", strings.ReplaceAll(name, "@", "/"))
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			info := &workflowInfo.Info{
				Capture:     true,
				Done:        true,
				SearchError: err.Error(),
			}

			*s.actionDetails.GetByID(name) = *info
			return
		}

		w.authorizeRequest(request, s) //TODO: this does not cause rash

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			info := &workflowInfo.Info{
				Capture:     true,
				Done:        true,
				SearchError: err.Error(),
			}

			*s.actionDetails.GetByID(name) = *info
			return
		}

		if response.StatusCode != 200 {
			info := &workflowInfo.Info{
				Capture: true,
				Done:    true,
			}

			switch response.StatusCode {
			case 404:
				info.SearchError = "Action not found"
			case 400:
				info.SearchError = "Invalid action name format"
			default:
				info.SearchError = fmt.Sprintf("Unexpected status code: %d", response.StatusCode)
			}

			*s.actionDetails.GetByID(name) = *info
			return
		}

		// 2. read response body
		output := make([]byte, 0)
		// read all content of response.Body
		// into output
		for {
			buffer := make([]byte, 1024)
			n, err := response.Body.Read(buffer)
			output = append(output, buffer[:n]...)
			if err != nil {
				break
			}
		}

		info := workflowInfo.Unmarshal(output)

		// 3. put it into s.actionDetails.GetByID(step.Uses)
		info.Capture = true
		info.Done = true
		// Reset this because most probably old options does not apply
		*s.actionDetails.GetByID(name) = *info
	}()
}

func (w *Widget) SearchActionBranches(name string, s *State) {
	go func() {
		s.branchesList = make([]string, 0)
		s.currentBranch = 0
		client := &http.Client{}
		url := fmt.Sprintf("https://api.github.com/repos/%s/branches", name)
		fmt.Println("URL", url)
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("crash: bad request")
			return
		}

		w.authorizeRequest(request)

		response, err := client.Do(request)
		if err != nil {
			fmt.Println("crash: can't do request")
			return
		}

		if response.StatusCode != 200 {
			fmt.Println("non-200 status code", response.StatusCode)
			return
		}

		output := make([]byte, 0)
		for {
			buffer := make([]byte, 1024)
			n, err := response.Body.Read(buffer)
			output = append(output, buffer[:n]...)
			if err != nil {
				break
			}
		}

		type branch struct {
			Name string `json:"name"`
		}

		// unmarshal output into branches
		// and put them into s.branchesList
		// and set s.currentBranch to 0
		// (if there are any branches)
		b := new([]branch)
		err = json.Unmarshal(output, b)
		if err != nil {
			fmt.Printf("Cannot unmarshal: %v\n", err)
			return
		}

		for _, branch := range *b {
			s.branchesList = append(s.branchesList, branch.Name)
		}
	}()
}

func (w *Widget) authorizeRequest(request *http.Request, state ...*State) bool {
	var s *State

	switch len(state) {
	case 0:
		s = w.GetState()
	case 1:
		s = state[0]
	default:
		panic("what the hell are you doing???")
	}

	// now set auth token
	if w.token == "" || s.InvalidToken {
		return false
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.token))
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	return true

}
