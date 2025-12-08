package window

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	schedule.AddSystem(sch, schedule.PreStartup, applyInitial).Label("window.applyInitial")
	schedule.AddSystem(sch, schedule.First, applyChanges).Label("window.applyChanges")
}
