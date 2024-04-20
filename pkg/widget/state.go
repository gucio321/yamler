package widget

import (
	"fmt"
	"github.com/gucio321/yamler/pkg/widget/workflowInfo"

	"github.com/AllenDang/giu"
	"github.com/gucio321/yamler/pkg/workflow"
)

type State struct {
	workflow      *workflow.Workflow
	code          string
	toggles       *SuperMap[bool]
	dropdowns     *SuperMap[int32]
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
		workflow:      w.w,
		toggles:       NewSuperMap[bool](),
		dropdowns:     NewSuperMap[int32](),
		actionDetails: NewSuperMap[workflowInfo.Info](),
		actionsWith:   NewSuperMap[string](),
	}

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
