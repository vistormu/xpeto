package pkg

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"

	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/hierarchy"
	"github.com/vistormu/xpeto/pkg/input"
	"github.com/vistormu/xpeto/pkg/sprite"
	"github.com/vistormu/xpeto/pkg/text"
)

func DefaultPkgs(w *ecs.World, sch *schedule.Scheduler) {
	// core
	asset.Pkg(w, sch)
	hierarchy.Pkg(w, sch)
	input.Pkg(w, sch)

	// dependent
	text.Pkg(w, sch)
	sprite.Pkg(w, sch)
}
