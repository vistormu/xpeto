package input

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/schedule"
)

func InputPlugin(ctx *core.Context, sch *schedule.Scheduler) {
	// resources
	core.AddResource(ctx, newKeyboard())
	core.AddResource(ctx, newMouse())

	// systems
	schedule.AddSystem(sch, schedule.PreUpdate, updateKeyboardInput)
	schedule.AddSystem(sch, schedule.PreUpdate, updateMouseInput)
	schedule.AddSystem(sch, schedule.PreUpdate, updateMouseCursor)
	schedule.AddSystem(sch, schedule.PreUpdate, updateMouseWheel)
}
