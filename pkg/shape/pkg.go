package shape

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/pkg/render"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// ellipse
	render.AddExtractionFn(w, extractEllipse)
	render.AddSortFn(w, sortEllipse)
	render.AddRenderFn(w, render.Opaque, drawEllipse)

	// path
	render.AddExtractionFn(w, extractPath)
	render.AddSortFn(w, sortPath)
	render.AddRenderFn(w, render.Opaque, drawPath)

	// rect
	render.AddExtractionFn(w, extractRect)
	render.AddSortFn(w, sortRect)
	render.AddRenderFn(w, render.Opaque, drawRect)
}
