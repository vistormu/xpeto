package game

import (
	"image/color"

	"github.com/vistormu/xpeto"
	"github.com/vistormu/xpeto/pkg/physics"
)

// ======
// config
// ======
type Config struct {
	PaddleSpeed  float64
	MinBallSpeed float64
	MaxBallSpeed float64
}

// ======
// paddle
// ======
type Paddle struct {
	IsLeft bool
}

func createLeftPaddle(w *xp.World) {
	e := xp.AddEntity(w)
	xp.AddComponent(w, e, Paddle{
		IsLeft: true,
	})
	xp.AddComponent(w, e, xp.NewRectShape(10, 30).
		AddFillSolid(color.RGBA{216, 166, 87, 255}).
		AddOrder(1, 0),
	)
	xp.AddComponent(w, e, xp.Transform{
		X: 20,
		Y: 50,
	})
	xp.AddComponent(w, e, physics.Velocity{})
	xp.AddComponent(w, e, physics.RigidBody{
		Type:        physics.Kinematic,
		Mass:        10,
		Restitution: 1,
		Friction:    0,
	})
}

func createRightPaddle(w *xp.World) {
	ww, wh := xp.GetVirtualWindowSize[float64](w)

	e := xp.AddEntity(w)
	xp.AddComponent(w, e, Paddle{
		IsLeft: false,
	})
	xp.AddComponent(w, e, xp.NewRectShape(10, 30).
		AddFillSolid(color.RGBA{125, 174, 163, 255}).
		AddOrder(1, 0),
	)
	xp.AddComponent(w, e, xp.Transform{
		X: ww - 20,
		Y: wh - 50,
	})
	xp.AddComponent(w, e, physics.Velocity{})
	xp.AddComponent(w, e, physics.RigidBody{
		Type:        physics.Kinematic,
		Mass:        10,
		Restitution: 1,
		Friction:    0,
	})
}

// ====
// ball
// ====
type Ball struct{}

func createBall(w *xp.World) {
	ww, wh := xp.GetVirtualWindowSize[float64](w)

	e := xp.AddEntity(w)
	xp.AddComponent(w, e, Ball{})
	xp.AddComponent(w, e, xp.Transform{
		X: ww / 2,
		Y: wh / 2,
	})
	xp.AddComponent(w, e, xp.NewCircleShape(5).
		AddFillSolid(color.RGBA{234, 105, 98, 255}).
		AddOrder(2, 0),
	)
	xp.AddComponent(w, e, physics.Velocity{})
	xp.AddComponent(w, e, physics.RigidBody{
		Type:        physics.Dynamic,
		Mass:        1,
		Restitution: 1,
		Friction:    0,
	})
	xp.AddComponent(w, e, xp.NewPathShape().
		AddStroke(color.White, 1).
		AddOrder(2, 0),
	)
}

// =====
// field
// =====

func createField(w *xp.World) {
	ww, wh := xp.GetVirtualWindowSize[float64](w)

	e := xp.AddEntity(w)
	xp.AddComponent(w, e, xp.NewRectShape(float32(ww), float32(wh)).
		AddFillSolid(color.RGBA{125, 174, 163, 255}).
		AddOrder(1, 0),
	)
	xp.AddComponent(w, e, xp.Transform{
		X: ww / 2,
		Y: wh / 2,
	})
}
