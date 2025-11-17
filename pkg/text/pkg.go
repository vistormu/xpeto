package text

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/pkg/render"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	render.AddExtractionFn(w, extractText)
	render.AddSortFn(w, sortText)
	render.AddRenderFn(w, render.Ui, drawText)
}
