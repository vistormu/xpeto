package time

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/schedule"
)

func TimePlugin(ctx *core.Context, sch *schedule.Scheduler) {
	// resources
	core.AddResource(ctx, new(Time))

	// systems
	clock := NewClock()
	schedule.AddSystem(sch, schedule.PreUpdate, clock.Update)
}
