package font

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/render"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	asset.AddLoaderFn(w, load, ".ttf")

	render.AddExtractionFn(w, extractText)
	render.AddSortFn(w, sortText)
	render.AddRenderFn(w, render.Ui, drawText)
}
