package render

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/schedule"
)

func RenderPlugin(ctx *core.Context, sch *schedule.Scheduler) {
	// resources
	core.AddResource(ctx, NewExtractor())

	// system
	schedule.AddSystem(sch, schedule.Draw, draw)
}
