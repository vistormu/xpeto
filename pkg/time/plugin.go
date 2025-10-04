package time

import (
	"time"

	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/schedule"
)

func TimePlugin(ctx *core.Context, sch *schedule.Scheduler) {
	// resources
	core.AddResource(ctx, &Time{FixedDelta: time.Second / 60})

	// systems
	clock := NewClock()
	schedule.AddSystem(sch, schedule.PreUpdate, clock.Update)
}
