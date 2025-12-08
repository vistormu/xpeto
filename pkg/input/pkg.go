package input

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	ecs.AddResource(w, newKeyboard())
	ecs.AddResource(w, newMouse())
}
