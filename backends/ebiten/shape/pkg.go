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

	// rect
	render.AddExtractionFn[ebiten.Image](w, extractRect)
	render.AddSortFn[ebiten.Image](w, sortRect)
	render.AddRenderFn(w, render.Opaque, drawRect)

	// path
	render.AddExtractionFn[ebiten.Image](w, extractPath)
	render.AddSortFn[ebiten.Image](w, sortPath)
	render.AddRenderFn(w, render.Opaque, drawPath)

	// line
	render.AddExtractionFn[ebiten.Image](w, extractLine)
	render.AddSortFn[ebiten.Image](w, sortLine)
	render.AddRenderFn(w, render.Opaque, drawLine)

	// segment
	render.AddExtractionFn[ebiten.Image](w, extractSegment)
	render.AddSortFn[ebiten.Image](w, sortSegment)
	render.AddRenderFn(w, render.Opaque, drawSegment)

	// arrow
	render.AddExtractionFn[ebiten.Image](w, extractArrow)
	render.AddSortFn[ebiten.Image](w, sortArrow)
	render.AddRenderFn(w, render.Opaque, drawArrow)
}
