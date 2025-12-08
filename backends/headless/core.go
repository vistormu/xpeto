package headless

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"

	"github.com/vistormu/xpeto/backends/headless/time"
)

func corePkgs(w *ecs.World, sch *schedule.Scheduler) {
	time.Pkg(w, sch)
}
