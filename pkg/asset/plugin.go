package asset

import (
	"github.com/vistormu/xpeto/internal/core"
)

type AssetPlugin struct{}

func (ap *AssetPlugin) Build(ctx *core.Context, sb *core.ScheduleBuilder) {
	core.AddResource(ctx, NewServer())
	sb.WithSystem("asset_loader", core.First, Update)
}
