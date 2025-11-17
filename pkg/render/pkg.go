package render

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// resources
	ecs.AddResource(w, newRenderer())

	// system
	schedule.AddSystem(sch, schedule.Draw, draw).Label("render.draw")
}
