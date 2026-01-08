package window

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// window
	ecs.AddResource(w, newRealWindow())
	ecs.AddResource(w, RealWindowObserved{})
	ecs.AddResource(w, newVirtualWindow())

	// scaling
	ecs.AddResource(w, newScaling())

	// viewport
	ecs.AddResource(w, ComputeViewport(w))
}
