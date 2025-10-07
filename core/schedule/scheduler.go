package schedule

import (
	"reflect"

	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"
)

type Scheduler struct {
	lastSchedule  *Schedule
	schedules     map[stage][]*Schedule
	stateMachines map[reflect.Type]any
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		lastSchedule:  nil,
		schedules:     make(map[stage][]*Schedule),
		stateMachines: make(map[reflect.Type]any),
	}
}

func (s *Scheduler) addSchedule(stage stage, schedule *Schedule) {
	_, ok := s.schedules[stage]
	if !ok {
		s.schedules[stage] = make([]*Schedule, 0)
	}

	s.schedules[stage] = append(s.schedules[stage], schedule)
	s.lastSchedule = schedule
}

func AddStateMachine[T comparable](sch *Scheduler, initial T) {
	sm := newStateMachine(initial)
	sch.stateMachines[reflect.TypeFor[T]()] = sm
	sch.addSchedule(stateTransition, &Schedule{
		Stage:  stateTransition,
		System: sm.update,
	})
	sch.addSchedule(postStartup, &Schedule{
		Stage:  postStartup,
		System: sm.startup,
	})
}

func getStateMachine[T comparable](sch *Scheduler) (*StateMachine[T], bool) {
	sm, ok := sch.stateMachines[reflect.TypeFor[T]()]
	if !ok {
		return nil, false
	}

	return sm.(*StateMachine[T]), true
}

func AddSystem(sch *Scheduler, stage Stage, system ecs.System) *Scheduler {
	st := stage(sch, system)
	if st != stateTransition {
		sch.addSchedule(st, &Schedule{
			Stage:  st,
			System: system,
		})
	}

	return sch
}

func (sch *Scheduler) RunIf(condition ConditionFn) *Scheduler {
	if sch.lastSchedule == nil {
		return sch
	}

	sch.lastSchedule.Condition = condition

	return sch
}

func (sch *Scheduler) Label(label string) *Scheduler {
	if sch.lastSchedule == nil {
		return sch
	}

	sch.lastSchedule.Label = label

	return sch
}

func (sch *Scheduler) After(label string) *Scheduler {
	if sch.lastSchedule == nil {
		// TODO: log something
		return sch
	}

	sch.lastSchedule.After = append(sch.lastSchedule.After, label)

	return sch
}

func (sch *Scheduler) Before(label string) *Scheduler {
	if sch.lastSchedule == nil {
		// TODO: log something
		return sch
	}

	sch.lastSchedule.Before = append(sch.lastSchedule.Before, label)

	return sch
}

func (s *Scheduler) run(ctx *core.Context, stages []stage) {
	for _, stage := range stages {
		// TODO: add sorting
		schedules := s.schedules[stage]
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

func (s *Scheduler) RunStartup(ctx *core.Context) {
	stages := []stage{preStartup, startup, postStartup}
	s.run(ctx, stages)
}

func (s *Scheduler) RunUpdate(ctx *core.Context) {
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
	s.run(ctx, stages)
}

func (s *Scheduler) RunDraw(ctx *core.Context) {
	stages := []stage{preDraw, draw, postDraw}
	s.run(ctx, stages)
}
