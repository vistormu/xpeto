package schedule

import (
	"github.com/vistormu/go-dsa/hashmap"
	"github.com/vistormu/xpeto/core/ecs"
)

// ========
// schedule
// ========
type Scheduler struct {
	lastSchedule *Schedule
	schedules    map[Stage][]*Schedule
	labelToId    map[string]uint64
	nextId       uint64
	extra        *hashmap.TypeMap
	fixedStepsFn func(*ecs.World) int
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		lastSchedule: nil,
		schedules:    make(map[Stage][]*Schedule),
		labelToId:    make(map[string]uint64),
		nextId:       1,
		extra:        hashmap.NewTypeMap(),
	}
}

func (sch *Scheduler) addSchedule(s *Schedule) {
	s.Id = sch.nextId
	sch.nextId++

	sch.lastSchedule = s

	if s.stage == empty {
		return
	}

	_, ok := sch.schedules[s.stage]
	if !ok {
		sch.schedules[s.stage] = make([]*Schedule, 0)
	}

	sch.schedules[s.stage] = append(sch.schedules[s.stage], s)
}

func (sch *Scheduler) run(w *ecs.World, stages []Stage) {
	for _, stage := range stages {
		for _, sch := range sch.schedules[stage] {
			if sch == nil {
				continue
			}

			execute := true
			for _, c := range sch.Conditions {
				if c == nil {
					continue
				}

				execute = execute && c(w)
			}

			if execute {
				ecs.SetSystemId(w, sch.Id)
				sch.System(w)
			}
		}
	}
}

// THIS METHOD SHOULD NOT BE CALLED
//
// it should only be called by the `xp.Game` struct
func (sch *Scheduler) RunStartup(w *ecs.World) {
	stages := []Stage{preStartup, startup, postStartup}
	sch.run(w, stages)
}

// THIS METHOD SHOULD NOT BE CALLED
//
// it should only be called by the `xp.Game` struct
func (sch *Scheduler) RunUpdate(w *ecs.World) {
	// first pass
	stages := []Stage{
		first,
		preUpdate,
		stateTransition,
	}
	sch.run(w, stages)

	// fixed pass
	steps := 0
	if sch.fixedStepsFn != nil {
		n := sch.fixedStepsFn(w)
		if n > 0 {
			steps = n
		}
	}
	stages = []Stage{
		fixedFirst,
		fixedPreUpdate,
		fixedUpdate,
		fixedPostUpdate,
		fixedLast,
	}

	for range steps {
		sch.run(w, stages)
	}

	// last pass
	stages = []Stage{
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
	stages := []Stage{preDraw, draw, postDraw}
	sch.run(w, stages)
}

// THIS METHOD SHOULD NOT BE CALLED
//
// it sets the number of steps for the fixed stage, set by the `Time` pkg
func (sch *Scheduler) SetFixedStepsFn(fn func(*ecs.World) int) {
	sch.fixedStepsFn = fn
}

// ===
// API
// ===
func AddSystem(sch *Scheduler, stage StageFn, system ecs.System) *Scheduler {
	s := newSchedule()
	s.System = system
	s.stage = stage(sch, s)

	sch.addSchedule(s)

	return sch
}

func (sch *Scheduler) RunIf(condition ConditionFn) *Scheduler {
	if sch.lastSchedule == nil {
		return sch
	}

	sch.lastSchedule.Conditions = append(sch.lastSchedule.Conditions, condition)

	return sch
}

func (sch *Scheduler) Label(label string) *Scheduler {
	if sch.lastSchedule == nil {
		return sch
	}

	sch.labelToId[label] = sch.lastSchedule.Id

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

func AddExtra[T any](sch *Scheduler, v T) {
	hashmap.Add(sch.extra, v)
}

func GetExtra[T any](sch *Scheduler) (*T, bool) {
	return hashmap.Get[T](sch.extra)
}
