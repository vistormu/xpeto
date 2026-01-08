package clock

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// settings
	cs := newClockSettings()
	ecs.AddResource(w, cs)

	// clocks
	ecs.AddResource(w, newRealClock(cs.Now())())
	ecs.AddResource(w, newVirtualClock(cs.Now())())
	ecs.AddResource(w, newFixedClock())

	// systems
	schedule.SetFixedStepsFn(sch, func(w *ecs.World) int {
		fixed, ok := ecs.GetResource[FixedClock](w)
		if !ok {
			return 0
		}
		return fixed.Steps
	})

	schedule.AddSystem(sch, schedule.First, sanitizeClockSettings,
		schedule.SystemOpt.Label("clock.sanitizeClockSettings"),
	)

	schedule.AddSystem(sch, schedule.First, tick,
		schedule.SystemOpt.Label("clock.tick"),
	)
}
