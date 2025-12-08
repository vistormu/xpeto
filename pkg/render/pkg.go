package render

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func Pkg[C any](w *ecs.World, sch *schedule.Scheduler) {
	// resources
	ecs.AddResource(w, newRenderer[C]())

	// system
	schedule.AddSystem(sch, schedule.Draw, render[C]).Label("render.draw")
}
