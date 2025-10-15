package vector

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/pkg/render"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// extractors
	render.AddExtractionFn(w, render.Opaque, extractCircles)
	render.AddExtractionFn(w, render.Opaque, extractRects)
}
