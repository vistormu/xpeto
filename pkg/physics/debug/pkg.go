package debug

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/pkg/render"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// resources
	ecs.AddResource(w, defaultSettings())

	// extractor
	render.AddExtractionFn[*ebiten.Image](w, extractAabb)
	render.AddSortFn[*ebiten.Image](w, sortAabb)
	render.AddRenderFn(w, render.Ui, drawAabb)

	// render.AddExtractionFn(w, render.Opaque, extractVelocities)
	// render.AddExtractionFn(w, render.Opaque, extractContacts)
	// render.AddExtractionFn(w, render.Opaque, extractGrid)
}
