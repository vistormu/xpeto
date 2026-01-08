package asset

import (
	"github.com/vistormu/xpeto/assets"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// resources
	ecs.AddResource(w, newServer())
	ecs.AddResource(w, newLoader())

	// default assets
	AddStaticFS(w, "default", assets.DefaultFS)

	// systems
	schedule.AddSystem(sch, schedule.First, readRequests,
		schedule.SystemOpt.Label("asset.readRequests"),
	)
	schedule.AddSystem(sch, schedule.First, loadResults,
		schedule.SystemOpt.Label("asset.loadResults"),
	)
}
