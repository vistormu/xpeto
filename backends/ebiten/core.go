package ebiten

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"

	"github.com/vistormu/xpeto/backends/ebiten/time"
	"github.com/vistormu/xpeto/backends/ebiten/window"
)

func corePkgs(w *ecs.World, sch *schedule.Scheduler) {
	time.Pkg(w, sch)
	window.Pkg(w, sch)
}
