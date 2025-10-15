package input

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// resources
	ecs.AddResource(w, newKeyboard())
	ecs.AddResource(w, newMouse())

	// systems
	schedule.AddSystem(sch, schedule.PreUpdate, updateKeyboardInput)
	schedule.AddSystem(sch, schedule.PreUpdate, updateMouseInput)
	schedule.AddSystem(sch, schedule.PreUpdate, updateMouseCursor)
	schedule.AddSystem(sch, schedule.PreUpdate, updateMouseWheel)
}
