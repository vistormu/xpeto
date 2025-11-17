package asset

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

// server is not concurrent safe, so the functions
// `AddFileSystem`... should be only called during game initialization
func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// resources
	ecs.AddResource(w, newServer())
	ecs.AddResource(w, newLoader())

	// systems
	schedule.AddSystem(sch, schedule.First, readRequests).Label("asset.readRequests")
	schedule.AddSystem(sch, schedule.First, loadResults).Label("asset.loadResults")
}
