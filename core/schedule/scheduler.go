package schedule

import (
	"github.com/vistormu/go-dsa/hashmap"
	"github.com/vistormu/xpeto/core/ecs"
)

type Scheduler struct {
	lastSchedule  *schedule
	schedules     map[stage][]*schedule
	stateMachines *hashmap.TypeMap
	labelToId     map[string]uint64
	nextId        uint64
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		lastSchedule:  nil,
		schedules:     make(map[stage][]*schedule),
		stateMachines: hashmap.NewTypeMap(),
		labelToId:     make(map[string]uint64),
		nextId:        1,
	}
}

func (sch *Scheduler) addSchedule(s *schedule) {
	_, ok := sch.schedules[s.stage]
	if !ok {
		sch.schedules[s.stage] = make([]*schedule, 0)
	}

	sch.schedules[s.stage] = append(sch.schedules[s.stage], s)
	sch.lastSchedule = s
}

func AddSystem(sch *Scheduler, stage StageFn, system ecs.System) *Scheduler {
	s := newSchedule()
	s.system = system
	s.id = sch.nextId
	s.stage = stage(sch, s)

	sch.nextId++

	if s.stage != stateTransition {
		sch.addSchedule(s)
	}

	return sch
}

func (sch *Scheduler) RunIf(condition ConditionFn) *Scheduler {
	if sch.lastSchedule == nil {
		return sch
	}

	sch.lastSchedule.conditions = append(sch.lastSchedule.conditions, condition)

	return sch
}

func (sch *Scheduler) Label(label string) *Scheduler {
	if sch.lastSchedule == nil {
		return sch
	}

	sch.labelToId[label] = sch.lastSchedule.id

	return sch
}

func (sch *Scheduler) After(label string) *Scheduler {
	if sch.lastSchedule == nil {
		return sch
	}

	id, ok := sch.labelToId[label]
	if !ok {
		return sch
	}

	sch.lastSchedule.after = append(sch.lastSchedule.after, id)

	return sch
}

func (sch *Scheduler) Before(label string) *Scheduler {
	if sch.lastSchedule == nil {
		return sch
	}

	id, ok := sch.labelToId[label]
	if !ok {
		return sch
	}

	sch.lastSchedule.before = append(sch.lastSchedule.before, id)

	return sch
}

func (sch *Scheduler) run(w *ecs.World, stages []stage) {
	for _, stage := range stages {
		for _, sch := range sch.schedules[stage] {
			if sch == nil {
				continue
			}

			execute := true
			for _, c := range sch.conditions {
				if c == nil {
					continue
				}

				execute = execute && c(w)
			}

			if execute {
				ecs.SetSystemId(w, sch.id)
				sch.system(w)
			}
		}
	}
}

// THIS METHOD SHOULD NOT BE CALLED
//
// it should only be called by the `xp.Game` struct
func (sch *Scheduler) RunStartup(w *ecs.World) {
	stages := []stage{preStartup, startup, postStartup}
	sch.run(w, stages)
}

// THIS METHOD SHOULD NOT BE CALLED
//
// it should only be called by the `xp.Game` struct
func (sch *Scheduler) RunUpdate(w *ecs.World) {
	stages := []stage{
		first,
		preUpdate,
		stateTransition,
		fixedFirst,
		fixedPreUpdate,
		fixedUpdate,
		fixedPostUpdate,
		fixedLast,
		update,
		postUpdate,
		last,
	}
	sch.run(w, stages)
}

// THIS METHOD SHOULD NOT BE CALLED
//
// it should only be called by the `xp.Game` struct
func (sch *Scheduler) RunDraw(w *ecs.World) {
	stages := []stage{preDraw, draw, postDraw}
	sch.run(w, stages)
}
