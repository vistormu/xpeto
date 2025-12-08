package font

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/render"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	asset.AddLoaderFn(w, load, ".ttf")

	render.AddExtractionFn[ebiten.Image](w, extractText)
	render.AddSortFn[ebiten.Image](w, sortText)
	render.AddRenderFn(w, render.Ui, drawText)
}
