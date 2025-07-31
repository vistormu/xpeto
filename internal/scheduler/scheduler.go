package scheduler

import (
	"github.com/vistormu/xpeto/internal/core"
)

type Scheduler struct {
	schedulesByStage map[core.Stage][]*Schedule
	stageOrder       []core.Stage
}

func NewScheduler(order []core.Stage) *Scheduler {
	return &Scheduler{
		schedulesByStage: make(map[core.Stage][]*Schedule),
		stageOrder:       order,
	}
}

func (s *Scheduler) WithSchedule(sch *Schedule) *Scheduler {
	_, ok := s.schedulesByStage[sch.Stage]
	if !ok {
		s.schedulesByStage[sch.Stage] = []*Schedule{}
	}

	s.schedulesByStage[sch.Stage] = append(s.schedulesByStage[sch.Stage], sch)

	return s
}

func (s *Scheduler) Run(ctx *core.Context) {
	for _, stage := range s.stageOrder {
		schedules := topoSort(s.schedulesByStage[stage])
		for _, sys := range schedules {
			if sys == nil {
				continue
			}
			if sys.Condition == nil || sys.Condition(ctx) {
				sys.System(ctx)
			}
		}
	}
}

func topoSort(schedules []*Schedule) []*Schedule {
	inDegree := map[string]int{}
	graph := map[string][]string{}
	index := map[string]*Schedule{}

	// build graph
	for _, s := range schedules {
		index[s.Name] = s
		for _, dep := range s.After {
			graph[dep] = append(graph[dep], s.Name)
			inDegree[s.Name]++
		}
		for _, dep := range s.Before {
			graph[s.Name] = append(graph[s.Name], dep)
			inDegree[dep]++
		}
	}

	queue := []string{}
	for _, s := range schedules {
		if inDegree[s.Name] == 0 {
			queue = append(queue, s.Name)
		}
	}

	order := []*Schedule{}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		order = append(order, index[n])
		for _, m := range graph[n] {
			inDegree[m]--
			if inDegree[m] == 0 {
				queue = append(queue, m)
			}
		}
	}

	// fall back to insertion order on cycles (dev error)
	if len(order) != len(schedules) {
		return schedules
	}
	return order
}
