package physics

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"
	"github.com/vistormu/xpeto/pkg/time"
	"github.com/vistormu/xpeto/pkg/transform"
)

func applyGravity(ctx *core.Context) {
	ps := core.MustResource[*PhysicsSettings](ctx)
	w := core.MustResource[*ecs.World](ctx)
	t := core.MustResource[*time.Time](ctx)

	entities := w.Query(ecs.And(
		ecs.Has[*Velocity](),
		ecs.Has[*RigidBody](),
	))

	for _, e := range entities {
		v, _ := ecs.GetComponent[*Velocity](w, e)
		rb, _ := ecs.GetComponent[*RigidBody](w, e)
		gs, ok := ecs.GetComponent[*GravityScale](w, e)

		if rb.Type != Dynamic {
			continue
		}

		scale := 1.0
		if ok {
			scale = gs.Value
		}

		v.Y += float64(ps.Gravity.Y) * scale * t.FixedDelta.Seconds()
		v.X += float64(ps.Gravity.X) * scale * t.FixedDelta.Seconds()
	}
}

func integrateVelocities(ctx *core.Context) {
	w := core.MustResource[*ecs.World](ctx)
	t := core.MustResource[*time.Time](ctx)

	entities := w.Query(ecs.And(
		ecs.Has[*Velocity](),
		ecs.Has[*RigidBody](),
		ecs.Has[*transform.Transform](),
	))

	for _, e := range entities {
		v, _ := ecs.GetComponent[*Velocity](w, e)
		rb, _ := ecs.GetComponent[*RigidBody](w, e)
		tr, _ := ecs.GetComponent[*transform.Transform](w, e)

		if rb.Type == Static {
			continue
		}

		tr.Position.X += float32(v.X) * float32(t.FixedDelta.Seconds())
		tr.Position.Y += float32(v.Y) * float32(t.FixedDelta.Seconds())
	}
}

func buildBroadPhase(ctx *core.Context) {

}

func narrowPhaseAABB(ctx *core.Context) {

}

func resolveContactsAABB(ctx *core.Context) {

}
