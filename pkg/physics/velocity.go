package physics

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/pkg/time"
	"github.com/vistormu/xpeto/core/pkg/transform"
)

type Velocity struct {
	X float64
	Y float64
}

func integrateVelocities(w *ecs.World) {
	clk, _ := ecs.GetResource[time.FixedClock](w)

	q := ecs.NewQuery3[Velocity, RigidBody, transform.Transform](w)

	for _, b := range q.Iter() {
		v := b.A()
		rb := b.B()
		tr := b.C()

		if rb.Type == Static {
			continue
		}

		tr.X += v.X * clk.Delta.Seconds()
		tr.Y += v.Y * clk.Delta.Seconds()
	}
}
