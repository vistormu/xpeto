package schedule

import (
	"github.com/vistormu/xpeto/internal/core"
)

func InState[T comparable](s T) ConditionFn {
	return func(ctx *core.Context) bool {
		current, ok := core.GetResource[*State[T]](ctx)
		if !ok {
			return false
		}
		return current.Get() == s
	}
}
