package vector

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/schedule"
	"github.com/vistormu/xpeto/pkg/render"
)

func VectorPlugin(ctx *core.Context, sch *schedule.Scheduler) {
	e, ok := core.GetResource[*render.Extractor](ctx)
	if !ok {
		panic("cannot use sprites without the render plugin")
	}

	// TODO: implement this better
	e.AddExtractionFn(render.Opaque, extractCircles)
	e.AddExtractionFn(render.Opaque, extractRects)
}
