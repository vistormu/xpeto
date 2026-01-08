package ebiten

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/pkg/render"

	"github.com/vistormu/xpeto/backends/ebiten/font"
	"github.com/vistormu/xpeto/backends/ebiten/image"
	"github.com/vistormu/xpeto/backends/ebiten/input"
	"github.com/vistormu/xpeto/backends/ebiten/shape"
)

func DefaultPkgs(w *ecs.World, sch *schedule.Scheduler) {
	render.Pkg[ebiten.Image](w, sch)
	input.Pkg(w, sch)
	image.Pkg(w, sch)
	font.Pkg(w, sch)
	shape.Pkg(w, sch)
}
