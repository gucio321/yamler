package widget

import (
	"fmt"
	"github.com/AllenDang/giu"
	"sort"
)

func (w *Widget) needs(jobName string) giu.Widget {
	s := w.GetState()
	jobNames := make([]string, 0)
	for name := range s.workflow.Jobs {
		if name != jobName {
			jobNames = append(jobNames, name)
		}
	}

	sort.Strings(jobNames)

	availableJobNames := make([]string, len(jobNames))
	copy(availableJobNames, jobNames)

	needs := s.workflow.Jobs[jobName].Needs
	ptrs := make([]*int32, len(needs))
	for i, need := range needs {
		ptrs[i] = new(int32)
		*ptrs[i] = -1
		for _, name := range jobNames {
			if name == need {
				*ptrs[i] = int32(i)
				availableJobNames[i] = ""
			}
		}
	}

	for i := range availableJobNames {
		if availableJobNames[i] == "" {
			availableJobNames = append(availableJobNames[:i], availableJobNames[i+1:]...)
		}
	}

	combos := make([]giu.Widget, len(needs))
	for i := range needs {
		i := i
		combos = append(combos, giu.Layout{
			giu.Combo(fmt.Sprintf("##%s%d", jobName, i),
				func() string {
					if *ptrs[i] == -1 {
						return "--"
					}

					return jobNames[*ptrs[i]]
				}(), jobNames, ptrs[i]),
		})

	}

	return giu.Layout{
		giu.Layout(combos),
		giu.Custom(func() {
			tmp := int32(-1)
			giu.Combo(fmt.Sprintf("##%sNew", jobName),
				func() string {
					if tmp == -1 {
						return "--"
					}

					return availableJobNames[tmp]
				}(), availableJobNames, &tmp).OnChange(func() {
				if tmp == -1 {
					return
				}

				s.workflow.Jobs[jobName].Needs = append(s.workflow.Jobs[jobName].Needs, availableJobNames[tmp])
			}).Build()
		}),
	}
}
