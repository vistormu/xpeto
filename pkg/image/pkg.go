package image

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/pkg/asset"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	asset.AddLoaderFn(w, loadImage, ".png", ".jpg", ".jpeg")
}
