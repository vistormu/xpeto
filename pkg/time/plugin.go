package time

import (
	"github.com/vistormu/xpeto/internal/core"
)

type TimePlugin struct{}

func (tp *TimePlugin) Build(ctx *core.Context, sb *core.ScheduleBuilder) {
	core.AddResource(ctx, new(Time))

	clock := NewClock()
	sb.WithSystem("time", core.PreUpdate, clock.Update)
}
