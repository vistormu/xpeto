package image

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/pkg/asset"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// loader
	asset.AddAssetLoader(w, ".png", loadImage)
	asset.AddAssetLoader(w, ".jpg", loadImage)
	asset.AddAssetLoader(w, ".jpeg", loadImage)
}
