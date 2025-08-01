package image

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/pkg/asset"
)

func ImagePlugin(ctx *core.Context, sb *core.ScheduleBuilder) {
	// loader
	as, ok := core.GetResource[*asset.Server](ctx)
	if !ok {
		return
	}

	as.AddLoader(".png", LoadImage)
	as.AddLoader(".jpeg", LoadImage)
	as.AddLoader(".jpg", LoadImage)

	// system
	renderer := NewRenderer()
	sb.NewSchedule().
		WithSystem("renderer_update", core.PostUpdate, renderer.Update)

	sb.NewSchedule().
		WithSystem("renderer_draw", core.Draw, renderer.Draw)
}
