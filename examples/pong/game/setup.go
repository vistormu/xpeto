package game

import (
	"github.com/vistormu/xpeto"
	"github.com/vistormu/xpeto/pkg/physics"
)

func setupPhysics(w *xp.World) {
	ps, _ := xp.GetResource[physics.Space](w)
	ps.CellHeight = 32
	ps.CellWidth = 32
	ps.Height = 300
	ps.Width = 400
}

func setupConfig(w *xp.World) {
	xp.AddResource(w, Config{
		PaddleSpeed:  200,
		MinBallSpeed: 200,
		MaxBallSpeed: 300,
	})
}

func setupField(w *xp.World) {
	createLeftPaddle(w)
	createRightPaddle(w)
	createBall(w)
}

func setupMiniPkg(w *xp.World, _ *xp.Scheduler) {
	xp.SetRealWindowSize(w, 800, 600)
	xp.SetVirtualWindowSize(w, 400, 300)

	setupPhysics(w)
	setupConfig(w)
	setupField(w)
}
