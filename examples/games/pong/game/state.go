package game

import (
	"github.com/vistormu/xpeto"
)

type gameState uint8

const (
	stateInitial gameState = iota
	stateGameOver
	statePlaying
)

func stateManager(w *xp.World) {
	curr, _ := xp.GetState[gameState](w)

	if curr == stateInitial {
		xp.SetNextState(w, stateGameOver)
		return
	}

	_, ok := xp.GetEvents[StartIntent](w)
	if ok && curr == stateGameOver {
		xp.SetNextState(w, statePlaying)
		return
	}

	_, ok = xp.GetEvents[ScoreEvent](w)
	if ok {
		xp.SetNextState(w, stateGameOver)
	}
}

func stateMiniPkg(_ *xp.World, sch *xp.Scheduler) {
	xp.AddStateMachine(sch, stateInitial)

	xp.AddSystem(sch, xp.PreUpdate, stateManager)
}
