package asset

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

// TODO: should i use get state?
func IsAssetLoaded[B any]() schedule.ConditionFn {
	return func(w *ecs.World) bool {
		_, ok := ecs.GetResource[B](w)
		return ok
	}
}
