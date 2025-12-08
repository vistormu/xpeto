package window

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// resources
	ecs.AddResource(w, RealWindow{
		Width:               1920,
		Height:              1080,
		FullScreen:          false,
		AntiAliasing:        false,
		VSync:               false,
		RunnableOnUnfocused: true,
	})
	ecs.AddResource(w, VirtualWindow{
		Width:  192,
		Height: 108,
	})
}
