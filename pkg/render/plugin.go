package render

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/schedule"
)

func RenderPlugin(ctx *core.Context, sch *schedule.Scheduler) {
	// system
	renderer := NewRenderer()
	schedule.AddSystem(sch, schedule.PostUpdate, renderer.Update)
	schedule.AddSystem(sch, schedule.Draw, renderer.Draw)
}
