package asset

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/schedule"
)

// TODO: should i use get state?
func IsLoaded[B any]() schedule.ConditionFn {
	return func(ctx *core.Context) bool {
		_, ok := core.GetResource[B](ctx)
		return ok
	}
}
