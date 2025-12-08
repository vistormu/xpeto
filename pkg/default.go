package pkg

import (
	"github.com/vistormu/xpeto/assets"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"

	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/hierarchy"
	"github.com/vistormu/xpeto/pkg/input"
)

func DefaultPkgs(w *ecs.World, sch *schedule.Scheduler) {
	asset.Pkg(w, sch)
	hierarchy.Pkg(w, sch)
	input.Pkg(w, sch)

	asset.AddStaticFS(w, "default", assets.DefaultFS)
}
