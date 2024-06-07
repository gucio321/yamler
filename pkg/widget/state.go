package widget

import (
	"encoding/json"
	"fmt"
	"github.com/AllenDang/giu"
	"github.com/gucio321/yamler/pkg/widget/workflowInfo"
	"github.com/gucio321/yamler/pkg/workflow"
	"net/http"
	"time"
)

type State struct {
	APILimits     *APILimits
	apiTimer      chan bool
	signature     bool
	workflow      *workflow.Workflow
	code          string
	toggles       *SuperMap[bool]
	permissions   *SuperMap[int32]
	dropdowns     *SuperMap[int32]
	currentBranch int32
	branchesList  []string
	actionDetails *SuperMap[workflowInfo.Info]
	actionsWith   *SuperMap[string]
	newJobName    string
}

func (s *State) Dispose() {
	// noting to do here
}

func (w *Widget) GetState() *State {
	if s := giu.GetState[State](giu.Context, w.stateID()); s != nil {
		return s
	}

	newState := w.newState()
	giu.SetState(giu.Context, w.stateID(), newState)

	return w.GetState()
}

func (w *Widget) newState() *State {
	s := &State{
		APILimits:     NewAPILimits(),
		apiTimer:      make(chan bool, 1),
		signature:     true,
		workflow:      w.w,
		toggles:       NewSuperMap[bool](),
		permissions:   NewSuperMap[int32](),
		dropdowns:     NewSuperMap[int32](),
		actionDetails: NewSuperMap[workflowInfo.Info](),
		actionsWith:   NewSuperMap[string](),
	}

	for jobID, job := range s.workflow.Jobs {
		runsOnIdxPtr := s.dropdowns.GetByID(w.jobRunsOnID(jobID))
		for osIdx, os := range GetOSs() {
			if job.RunsOn == workflow.OS(os) {
				*runsOnIdxPtr = int32(osIdx)
				break
			}
		}
		for stepIdx, step := range job.Steps {
			if step.Uses != "" {
				w.SearchActionInputs(step.Uses, s)
				// also, fill  with's
				for key, value := range step.With {
					jobID := jobID
					stepIdx := stepIdx
					*s.actionsWith.GetByID(fmt.Sprintf("%s%s%d%s%s", key, step.Uses, stepIdx, jobID, w.id)) = value
				}
			}
		}
	}

	if err := w.updateRequestsLimitNoStateRace(s); err != nil {
		fmt.Println("Failed to get requests limit:", err)
	}

	go func() {
		for {
			select {
			case <-time.NewTimer(1 * time.Minute).C:
				if err := w.updateRequestsLimit(); err != nil {
					fmt.Println("Failed to get requests limit:", err)
				}
			case <-s.apiTimer:
				return
			}
		}
	}()

	return s
}

func (w *Widget) stateID() string {
	return fmt.Sprintf("%s_state", w.id)
}

type SuperMap[T any] map[string]*T

func NewSuperMap[T any]() *SuperMap[T] {
	m := make(SuperMap[T])
	return &m
}

func (s *SuperMap[T]) GetByID(id string) *T {
	if v, ok := (*s)[id]; ok {
		return v
	}

	newV := new(T)
	(*s)[id] = newV

	return newV
}

type APILimits struct {
	Limit     int
	Remaining int
	Reset     int
}

func NewAPILimits() *APILimits {
	return &APILimits{}
}

func (a *APILimits) Dec() {
	a.Remaining--
}

func (w *Widget) updateRequestsLimit() error {
	s := w.GetState()
	return w.updateRequestsLimitNoStateRace(s)
}

func (w *Widget) updateRequestsLimitNoStateRace(s *State) error {
	// talk to GItHub api and get current requests limit
	url := "https://api.github.com/rate_limit"
	type response struct {
		Resources struct {
			Core struct {
				Limit     int `json:"limit"`
				Remaining int `json:"remaining"`
				Reset     int `json:"reset"`
			} `json:"core"`
		} `json:"resources"`
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	w.authorizeRequest(request)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return err
	}

	s.APILimits.Limit = r.Resources.Core.Limit
	s.APILimits.Remaining = r.Resources.Core.Remaining
	s.APILimits.Reset = r.Resources.Core.Reset

	return nil
}
