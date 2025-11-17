package game

import (
	"github.com/vistormu/xpeto"
	"github.com/vistormu/xpeto/pkg/physics"
	// "github.com/vistormu/xpeto/pkg/physics/debug"
)

func Pkg(w *xp.World, sch *xp.Scheduler) {
	physics.Pkg(w, sch)
	// debug.Pkg(w, sch)

	stateMiniPkg(w, sch)
	uiMiniPkg(w, sch)
	inputMiniPkg(w, sch)
	setupMiniPkg(w, sch)
	scoreMiniPkg(w, sch)
	movementMiniPkg(w, sch)
}
