package schedule

import (
	"github.com/vistormu/xpeto/core/ecs"
)

type ConditionFn = func(*ecs.World) bool

func Once() ConditionFn {
	done := false
	return func(w *ecs.World) bool {
		if done {
			return false
		}
		done = true
		return true
	}
}

func OnceWhen(fn ConditionFn) ConditionFn {
	done := false
	return func(w *ecs.World) bool {
		if done {
			return false
		}
		if fn(w) {
			done = true
			return true
		}
		return false
	}
}
