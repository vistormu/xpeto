package game

import (
	"pong/assets"

	"github.com/vistormu/xpeto"
)

func Pkg(w *xp.World, sch *xp.Scheduler) {
	// resources

	// state machines
	xp.AddStateMachine(sch, initial)

	xp.AddFileSystem(w, assets.Assets)

	// systems
	xp.AddSystem(sch, xp.Startup, loadASsets)
	xp.AddSystem(sch, xp.Startup, initPhysiscs)
	xp.AddSystem(sch, xp.Startup, initWindow)
	xp.AddSystem(sch, xp.Update, createInitialGameplayScene).
		RunIf(xp.OnceWhen(xp.IsAssetLoaded[Fonts]()))
	xp.AddSystem(sch, xp.Update, stateManager)
	xp.AddSystem(sch, xp.OnEnter(playing), onEnterPlaying)
	xp.AddSystem(sch, xp.Update, gameplay)
}
