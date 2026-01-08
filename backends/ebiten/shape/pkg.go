package shape

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/pkg/render"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	render.AddExtractionFn[ebiten.Image](w, extractShape)
	render.AddSortFn[ebiten.Image](w, sortShape)
	render.AddRenderFn(w, render.Opaque, drawShape)
}
