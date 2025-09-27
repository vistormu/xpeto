package render

import (
	"github.com/vistormu/xpeto/internal/core"
)

func RenderPlugin(ctx *core.Context, sb *core.ScheduleBuilder) {
	// system
	renderer := NewRenderer()
	sb.NewSchedule().
		WithSystem("renderer_update", core.PostUpdate, renderer.Update)

	sb.NewSchedule().
		WithSystem("renderer_draw", core.Draw, renderer.Draw)
}
