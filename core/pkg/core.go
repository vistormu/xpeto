package pkg

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"

	"github.com/vistormu/xpeto/core/pkg/event"
	"github.com/vistormu/xpeto/core/pkg/log"
	"github.com/vistormu/xpeto/core/pkg/time"
	"github.com/vistormu/xpeto/core/pkg/window"
)

func CorePkgs(w *ecs.World, sch *schedule.Scheduler) {
	// core: no internal dependencies
	event.Pkg(w, sch)
	time.Pkg(w, sch)
	window.Pkg(w, sch)

	// features: depends on core packages
	log.Pkg(w, sch)
}
