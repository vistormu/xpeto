package window

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// window
	ecs.AddResource(w, RealWindow{
		Title:               "xpeto app",
		Width:               800,
		Height:              600,
		FullScreen:          false,
		AntiAliasing:        false,
		VSync:               false,
		RunnableOnUnfocused: true,
		ResizingMode:        ResizingModeDisabled,
		SizeLimits:          SizeLimits{-1, -1, -1, -1},
		Action:              ActionNone,
	})
	ecs.AddResource(w, RealWindowObserved{})
	ecs.AddResource(w, VirtualWindow{
		Width:  800,
		Height: 600,
	})

	// scaling
	ecs.AddResource(w, Scaling{
		Mode:       ScalingInteger,
		SnapPixels: true,
	})

	// viewport
	ecs.AddResource(w, ComputeViewport(w))
}
