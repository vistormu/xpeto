package physics

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/time"
)

type Gravity struct {
	X float64
	Y float64
}

type GravityScale struct {
	X float64
	Y float64
}

func applyGravity(w *ecs.World) {
	g, _ := ecs.GetResource[Gravity](w)
	clk, _ := ecs.GetResource[time.FixedClock](w)

	q := ecs.NewQuery2[Velocity, RigidBody](w)

	for _, b := range q.Iter() {
		v, rb := b.Components()

		if rb.Type != Dynamic {
			continue
		}

		gs, ok := ecs.GetComponent[GravityScale](w, b.Entity())
		if !ok {
			gs = &GravityScale{X: 1, Y: 1}
		}

		v.X += g.X * gs.X * clk.Delta.Seconds()
		v.Y += g.Y * gs.Y * clk.Delta.Seconds()
	}
}
