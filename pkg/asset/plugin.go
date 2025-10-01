package asset

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/schedule"
)

func AssetPlugin(ctx *core.Context, sch *schedule.Scheduler) {
	// resources
	core.AddResource(ctx, NewServer())

	// systems
	schedule.AddSystem(sch, schedule.First, Update)
}
