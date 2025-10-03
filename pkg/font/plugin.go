package font

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/schedule"
	"github.com/vistormu/xpeto/pkg/asset"
)

func FontPlugin(ctx *core.Context, sch *schedule.Scheduler) {
	// loader
	as, ok := core.GetResource[*asset.Server](ctx)
	if !ok {
		return
	}

	as.AddLoader(".ttf", loadText)
}
