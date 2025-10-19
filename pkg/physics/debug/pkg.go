package debug

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/pkg/render"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// resources
	ecs.AddResource(w, defaultSettings())

	// extractor
	render.AddExtractionFn(w, render.Opaque, extractAabb)
	render.AddExtractionFn(w, render.Opaque, extractVelocities)
	render.AddExtractionFn(w, render.Opaque, extractContacts)
	render.AddExtractionFn(w, render.Opaque, extractGrid)
}
