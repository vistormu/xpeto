package image

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/render"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	asset.AddLoaderFn(w, load, ".png", ".jpg", ".jpeg")

	render.AddExtractionFn(w, extractSprite)
	render.AddSortFn(w, sortSprite)
	render.AddRenderFn(w, render.Opaque, drawSprite)
}
