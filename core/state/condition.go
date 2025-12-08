package state

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func InState[T comparable](s T) schedule.ConditionFn {
	return func(w *ecs.World) bool {
		current, ok := GetState[T](w)
		if !ok {
			return false
		}
		return current == s
	}
}
