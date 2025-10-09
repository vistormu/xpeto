package state

import (
	"github.com/vistormu/xpeto/core/schedule"
)

func OnExit[T comparable](state T) schedule.StageFn {
	return func(sch *schedule.Scheduler, s *schedule.Schedule) schedule.Stage {
		sm, ok := schedule.GetExtra[stateMachine[T]](sch)

		if !ok {
			return schedule.Stage(0)
		}

		sm.addOnExit(state, s)

		return schedule.Stage(0)
	}
}
func OnTransition[T comparable](from, to T) schedule.StageFn {
	return func(sch *schedule.Scheduler, s *schedule.Schedule) schedule.Stage {
		sm, ok := schedule.GetExtra[stateMachine[T]](sch)
		if !ok {
			return schedule.Stage(0)
		}

		sm.addOnTransition(from, to, s)

		return schedule.Stage(0)
	}
}
func OnEnter[T comparable](state T) schedule.StageFn {
	return func(sch *schedule.Scheduler, s *schedule.Schedule) schedule.Stage {
		sm, ok := schedule.GetExtra[stateMachine[T]](sch)
		if !ok {
			return schedule.Stage(0)
		}

		sm.addOnEnter(state, s)

		return schedule.Stage(0)
	}
}
