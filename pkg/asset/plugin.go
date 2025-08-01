package asset

import (
	"github.com/vistormu/xpeto/internal/core"
)

func AssetPlugin(ctx *core.Context, sb *core.ScheduleBuilder) {
	// resources
	core.AddResource(ctx, NewServer())

	// systems
	sb.NewSchedule().
		WithSystem("asset_loader", core.First, Update)
}
