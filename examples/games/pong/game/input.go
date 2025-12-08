package game

import (
	"github.com/vistormu/xpeto"
)

type StartIntent struct{}
type PauseIntent struct{}
type MoveIntent struct {
	IsLeft bool
	IsUp   bool
}

func getInput(w *xp.World) {
	kb, _ := xp.GetResource[xp.Keyboard](w)

	if kb.IsJustPressed(xp.KeyEnter) {
		xp.AddEvent(w, StartIntent{})
	}

	if kb.IsPressed(xp.KeyW) {
		xp.AddEvent(w, MoveIntent{
			IsLeft: true,
			IsUp:   true,
		})
	}

	if kb.IsPressed(xp.KeyS) {
		xp.AddEvent(w, MoveIntent{
			IsLeft: true,
			IsUp:   false,
		})
	}

	if kb.IsPressed(xp.KeyArrowUp) {
		xp.AddEvent(w, MoveIntent{
			IsLeft: false,
			IsUp:   true,
		})
	}

	if kb.IsPressed(xp.KeyArrowDown) {
		xp.AddEvent(w, MoveIntent{
			IsLeft: false,
			IsUp:   false,
		})
	}
}

func inputMiniPkg(_ *xp.World, sch *xp.Scheduler) {
	xp.AddSystem(sch, xp.PreUpdate, getInput)
}
