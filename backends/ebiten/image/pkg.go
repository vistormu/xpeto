package image

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/render"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	asset.AddLoaderFn(w, load, ".png", ".jpg", ".jpeg")

	render.AddExtractionFn[ebiten.Image](w, extractSprite)
	render.AddSortFn[ebiten.Image](w, sortSprite)
	render.AddRenderFn(w, render.Opaque, drawSprite)
}
