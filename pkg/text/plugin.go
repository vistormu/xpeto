package text

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/pkg/asset"
)

func TextPlugin(ctx *core.Context, sb *core.ScheduleBuilder) {
	// loader
	as, ok := core.GetResource[*asset.Server](ctx)
	if !ok {
		return
	}

	as.AddLoader(".ttf", loadText)
}
