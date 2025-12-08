package shape

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/pkg/render"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// ellipse
	render.AddExtractionFn[ebiten.Image](w, extractEllipse)
	render.AddSortFn[ebiten.Image](w, sortEllipse)
	render.AddRenderFn(w, render.Opaque, drawEllipse)

	// path
	render.AddExtractionFn[ebiten.Image](w, extractPath)
	render.AddSortFn[ebiten.Image](w, sortPath)
	render.AddRenderFn(w, render.Opaque, drawPath)

	// rect
	render.AddExtractionFn[ebiten.Image](w, extractRect)
	render.AddSortFn[ebiten.Image](w, sortRect)
	render.AddRenderFn(w, render.Opaque, drawRect)
}
