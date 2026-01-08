package clock

import (
	"time"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func EveryNFrames(n uint64) schedule.ConditionFn {
	if n == 0 {
		n = 1
	}
	return func(w *ecs.World) bool {
		vc, ok := ecs.GetResource[VirtualClock](w)
		if !ok {
			return false
		}
		return vc.Frame%n == 0
	}
}

func EveryNFixedFrames(n uint64) schedule.ConditionFn {
	if n == 0 {
		n = 1
	}
	return func(w *ecs.World) bool {
		fc, ok := ecs.GetResource[FixedClock](w)
		if !ok {
			return false
		}
		return fc.Frame%n == 0
	}
}

func OnceAfterElapsed(d time.Duration) schedule.ConditionFn {
	return schedule.OnceWhen(AfterElapsed(d))
}

func AfterElapsed(d time.Duration) schedule.ConditionFn {
	if d < 0 {
		d = 0
	}
	return func(w *ecs.World) bool {
		vc, ok := ecs.GetResource[VirtualClock](w)
		if !ok {
			return false
		}
		return vc.Elapsed >= d
	}
}

func OnceAfterRealElapsed(d time.Duration) schedule.ConditionFn {
	return schedule.OnceWhen(AfterRealElapsed(d))
}

func AfterRealElapsed(d time.Duration) schedule.ConditionFn {
	if d < 0 {
		d = 0
	}
	return func(w *ecs.World) bool {
		rc, ok := ecs.GetResource[RealClock](w)
		if !ok {
			return false
		}
		return rc.Elapsed >= d
	}
}

func EveryDuration(d time.Duration) schedule.ConditionFn {
	if d <= 0 {
		d = time.Nanosecond
	}

	var next time.Duration
	var init bool

	return func(w *ecs.World) bool {
		vc, ok := ecs.GetResource[VirtualClock](w)
		if !ok {
			return false
		}

		if !init {
			next = d
			init = true
		}

		if vc.Elapsed >= next {
			k := (vc.Elapsed - next) / d
			next += (k + 1) * d
			return true
		}

		return false
	}
}

func EveryFixedSteps() schedule.ConditionFn {
	return func(w *ecs.World) bool {
		fc, ok := ecs.GetResource[FixedClock](w)
		if !ok {
			return false
		}
		return fc.Steps > 0
	}
}
