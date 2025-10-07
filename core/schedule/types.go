package schedule

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"
)

// ========
// schedule
// ========
type ConditionFn = func(*core.Context) bool

type Schedule struct {
	Label     string
	Stage     stage
	System    ecs.System
	Before    []string
	After     []string
	Condition ConditionFn
}

// =====
// stage
// =====
type Stage = func(*Scheduler, ecs.System) stage

// startup
func PreStartup(*Scheduler, ecs.System) stage  { return preStartup }
func Startup(*Scheduler, ecs.System) stage     { return startup }
func PostStartup(*Scheduler, ecs.System) stage { return postStartup }

// update
func First(*Scheduler, ecs.System) stage     { return first }
func PreUpdate(*Scheduler, ecs.System) stage { return preUpdate }

func OnExit[T comparable](state T) Stage {
	return func(sch *Scheduler, system ecs.System) stage {
		sm, ok := getStateMachine[T](sch)
		if !ok {
			return stateTransition
		}

		_, ok = sm.onExit[state]
		if !ok {
			sm.onExit[state] = make([]ecs.System, 0)
		}

		sm.onExit[state] = append(sm.onExit[state], system)

		return stateTransition
	}
}
func OnTransition[T comparable](from, to T) Stage {
	return func(sch *Scheduler, system ecs.System) stage {
		sm, ok := getStateMachine[T](sch)
		if !ok {
			return stateTransition
		}

		key := transition[T]{from: from, to: to}

		_, ok = sm.onTransition[key]
		if !ok {
			sm.onTransition[key] = make([]ecs.System, 0)
		}

		sm.onTransition[key] = append(sm.onTransition[key], system)

		return stateTransition
	}
}
func OnEnter[T comparable](state T) Stage {
	return func(sch *Scheduler, system ecs.System) stage {
		sm, ok := getStateMachine[T](sch)
		if !ok {
			return stateTransition
		}

		_, ok = sm.onEnter[state]
		if !ok {
			sm.onEnter[state] = make([]ecs.System, 0)
		}

		sm.onEnter[state] = append(sm.onEnter[state], system)

		return stateTransition
	}
}

func FixedFirst(*Scheduler, ecs.System) stage      { return fixedFirst }
func FixedPreUpdate(*Scheduler, ecs.System) stage  { return fixedPreUpdate }
func FixedUpdate(*Scheduler, ecs.System) stage     { return fixedUpdate }
func FixedPostUpdate(*Scheduler, ecs.System) stage { return fixedPostUpdate }
func FixedLast(*Scheduler, ecs.System) stage       { return fixedLast }

func Update(*Scheduler, ecs.System) stage     { return update }
func PostUpdate(*Scheduler, ecs.System) stage { return postUpdate }
func Last(*Scheduler, ecs.System) stage       { return last }

// draw
func PreDraw(*Scheduler, ecs.System) stage  { return preDraw }
func Draw(*Scheduler, ecs.System) stage     { return draw }
func PostDraw(*Scheduler, ecs.System) stage { return postDraw }

type stage uint32

const (
	preStartup stage = iota
	startup
	postStartup

	first
	preUpdate

	stateTransition

	fixedFirst
	fixedPreUpdate
	fixedUpdate
	fixedPostUpdate
	fixedLast

	update
	postUpdate
	last

	preDraw
	draw
	postDraw
)

// =====
// state
// =====
type State[T comparable] struct {
	current T
}

func (s *State[T]) Get() T {
	return s.current
}

type NextState[T comparable] struct {
	next    *T
	pending bool
}

func (ns *NextState[T]) Set(s T) {
	ns.next = &s
	ns.pending = true
}

type transition[T comparable] struct {
	from, to T
}
