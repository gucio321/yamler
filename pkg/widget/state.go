package widget

import (
	"fmt"

	"github.com/AllenDang/giu"
	"github.com/gucio321/yamler/pkg/workflow"
)

type State struct {
	workflow *workflow.Workflow
	code     string
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
		workflow: &workflow.Workflow{}, // TODO
	}

	return s
}

func (w *Widget) stateID() string {
	return fmt.Sprintf("%s_state", w.id)
}
