package input

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// resources
	ecs.AddResource(w, newKeyboard())
	ecs.AddResource(w, newMouse())
	ecs.AddResource(w, newGamepads())

	// pipeline
	schedule.AddSystem(sch, schedule.First, beginFrame,
		schedule.SystemOpt.Label("input.beginFrame"),
	)
	schedule.AddSystem(sch, schedule.PreUpdate, applyEvents,
		schedule.SystemOpt.Label("input.applyEvents"),
	)
	schedule.AddSystem(sch, schedule.PreUpdate, watchFocus,
		schedule.SystemOpt.Label("input.watchFocus"),
	)
	schedule.AddSystem(sch, schedule.PreUpdate, compute,
		schedule.SystemOpt.Label("input.compute"),
	)
	schedule.AddSystem(sch, schedule.Last, endFrame,
		schedule.SystemOpt.Label("input.endFrame"),
	)
}
