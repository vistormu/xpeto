package schedule

import (
	"github.com/vistormu/go-dsa/hashmap"
)

type StageFn = func(*Scheduler, *schedule) stage
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

// startup
func PreStartup(*Scheduler, *schedule) stage  { return preStartup }
func Startup(*Scheduler, *schedule) stage     { return startup }
func PostStartup(*Scheduler, *schedule) stage { return postStartup }

// update
func First(*Scheduler, *schedule) stage     { return first }
func PreUpdate(*Scheduler, *schedule) stage { return preUpdate }

func OnExit[T comparable](state T) StageFn {
	return func(sch *Scheduler, s *schedule) stage {
		sm, ok := hashmap.Get[stateMachine[T]](sch.stateMachines)
		if !ok {
			return stateTransition
		}

		sm.addOnExit(state, s)

		return stateTransition
	}
}
func OnTransition[T comparable](from, to T) StageFn {
	return func(sch *Scheduler, s *schedule) stage {
		sm, ok := hashmap.Get[stateMachine[T]](sch.stateMachines)
		if !ok {
			return stateTransition
		}

		sm.addOnTransition(from, to, s)

		return stateTransition
	}
}
func OnEnter[T comparable](state T) StageFn {
	return func(sch *Scheduler, s *schedule) stage {
		sm, ok := hashmap.Get[stateMachine[T]](sch.stateMachines)
		if !ok {
			return stateTransition
		}

		sm.addOnEnter(state, s)

		return stateTransition
	}
}

func FixedFirst(*Scheduler, *schedule) stage      { return fixedFirst }
func FixedPreUpdate(*Scheduler, *schedule) stage  { return fixedPreUpdate }
func FixedUpdate(*Scheduler, *schedule) stage     { return fixedUpdate }
func FixedPostUpdate(*Scheduler, *schedule) stage { return fixedPostUpdate }
func FixedLast(*Scheduler, *schedule) stage       { return fixedLast }

func Update(*Scheduler, *schedule) stage     { return update }
func PostUpdate(*Scheduler, *schedule) stage { return postUpdate }
func Last(*Scheduler, *schedule) stage       { return last }

// draw
func PreDraw(*Scheduler, *schedule) stage  { return preDraw }
func Draw(*Scheduler, *schedule) stage     { return draw }
func PostDraw(*Scheduler, *schedule) stage { return postDraw }
