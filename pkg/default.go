package pkg

import (
	"github.com/vistormu/xpeto/assets"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"

	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/font"
	"github.com/vistormu/xpeto/pkg/hierarchy"
	"github.com/vistormu/xpeto/pkg/image"
	"github.com/vistormu/xpeto/pkg/input"
	"github.com/vistormu/xpeto/pkg/render"
	"github.com/vistormu/xpeto/pkg/shape"
)

func DefaultPkgs(w *ecs.World, sch *schedule.Scheduler) {
	// core: no dependencies
	asset.Pkg(w, sch)
	render.Pkg(w, sch)
	input.Pkg(w, sch)
	hierarchy.Pkg(w, sch)

	// semi-core: depends only on core
	image.Pkg(w, sch)
	font.Pkg(w, sch)
	shape.Pkg(w, sch)

	asset.AddStaticFS(w, "default", assets.DefaultFS)
}
