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
	schedule.AddSystem(sch, schedule.First, updateKeyboardInput).Label("input.updateKeyboardInput")
	schedule.AddSystem(sch, schedule.First, updateMouseInput).Label("input.updateMouseInput")
	schedule.AddSystem(sch, schedule.First, updateMouseCursor).Label("input.updateMouseCursor")
	schedule.AddSystem(sch, schedule.First, updateMouseWheel).Label("input.updateMouseWheel")
}
