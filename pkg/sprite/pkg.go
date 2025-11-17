package sprite

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/pkg/render"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// extraction
	render.AddExtractionFn(w, extractSprite)
	render.AddSortFn(w, sortSprite)
	render.AddRenderFn(w, render.Opaque, drawSprite)
}
