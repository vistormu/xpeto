package font

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/pkg/asset"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// loader
	asset.AddAssetLoader(w, ".ttf", loadFont)
}
