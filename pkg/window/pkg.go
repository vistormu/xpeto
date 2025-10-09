package window

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// resources
	ecs.AddResource(w, WindowSettings{
		Width:        1920,
		Height:       1080,
		VWidth:       192,
		VHeight:      108,
		FullScreen:   false,
		AntiAliasing: false,
	})

	// systems
	schedule.AddSystem(sch, schedule.PreStartup, applyInitialSettings)
	schedule.AddSystem(sch, schedule.PreUpdate, applyChanges)
}
