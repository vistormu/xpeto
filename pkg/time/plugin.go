package time

import (
	"github.com/vistormu/xpeto/internal/core"
)

func TimePlugin(ctx *core.Context, sb *core.ScheduleBuilder) {
	// resources
	core.AddResource(ctx, new(Time))

	// systems
	clock := NewClock()
	sb.NewSchedule().
		WithSystem("time", core.PreUpdate, clock.Update)
}
