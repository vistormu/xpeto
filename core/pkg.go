package core

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/core/window"

	"github.com/vistormu/xpeto/core/event"
	"github.com/vistormu/xpeto/core/log"
	"github.com/vistormu/xpeto/core/time"
)

type Pkg = func(*ecs.World, *schedule.Scheduler)

func CorePkgs(w *ecs.World, sch *schedule.Scheduler) {
	// core: no internal dependencies
	event.Pkg(w, sch)
	time.Pkg(w, sch)
	window.Pkg(w, sch)

	// features: depends on core packages
	log.Pkg(w, sch)
}
