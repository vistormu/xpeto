package image

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/schedule"
	"github.com/vistormu/xpeto/pkg/asset"
)

func ImagePlugin(ctx *core.Context, sb *schedule.Scheduler) {
	// loader
	as, ok := core.GetResource[*asset.Server](ctx)
	if !ok {
		return
	}

	as.AddLoader(".png", loadImage)
	as.AddLoader(".jpeg", loadImage)
	as.AddLoader(".jpg", loadImage)
}
