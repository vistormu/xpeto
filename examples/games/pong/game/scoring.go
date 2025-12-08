package game

import (
	"github.com/vistormu/xpeto"
)

type Score struct {
	Left  int
	Right int
}

type ScoreEvent struct{}

func trackScore(w *xp.World) {

}

func scoreMiniPkg(w *xp.World, sch *xp.Scheduler) {
	xp.AddResource(w, Score{})

	xp.AddSystem(sch, xp.Update, trackScore)
}
