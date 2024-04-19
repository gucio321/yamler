package widget

import (
	"fmt"

	"github.com/AllenDang/giu"
	"github.com/gucio321/yamler/pkg/workflow"
)

type State struct {
	workflow *workflow.Workflow
	code     string
	toggles  *SuperMap
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
		toggles:  NewSuperMap(),
	}

	return s
}

func (w *Widget) stateID() string {
	return fmt.Sprintf("%s_state", w.id)
}

type SuperMap map[string]*bool

func NewSuperMap() *SuperMap {
	m := SuperMap(make(map[string]*bool))
	return &m
}

func (s *SuperMap) GetByID(id string) *bool {
	if v, ok := (*s)[id]; ok {
		return v
	}

	newV := new(bool)
	(*s)[id] = newV

	return newV
}
