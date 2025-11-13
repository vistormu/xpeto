package time

import (
	"time"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// resources
	ecs.AddResource(w, ClockSettings{
		FixedDelta:  time.Second / 60,
		Scale:       1.0,
		Paused:      false,
		SyncWithFps: false,
		MaxDelta:    time.Millisecond * 100,
	})
	ecs.AddResource(w, RealClock{})
	ecs.AddResource(w, VirtualClock{})
	ecs.AddResource(w, FixedClock{
		Delta:       time.Second / 60,
		Accumulator: 0,
		MaxSteps:    8,
	})

	sch.SetFixedStepsFn(func(w *ecs.World) int {
		fixed, _ := ecs.GetResource[FixedClock](w)
		return fixed.Steps
	})

	// systems
	schedule.AddSystem(sch, schedule.PreStartup, applyInitialSettings)
	schedule.AddSystem(sch, schedule.PreUpdate, applyChanges)
	schedule.AddSystem(sch, schedule.First, tick)
}
