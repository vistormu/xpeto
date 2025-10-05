package physics

import (
	"math"

	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"
	"github.com/vistormu/xpeto/pkg/time"
	"github.com/vistormu/xpeto/pkg/transform"
)

func applyGravity(ctx *core.Context) {
	ps := core.MustResource[*Settings](ctx)
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

type contactPairs struct {
	a, b       ecs.Entity
	penX, penY float64
}

type CollisionSolver struct {
	cells    map[[2]int][]ecs.Entity
	contacts []contactPairs
}

func NewCollisionSolver() *CollisionSolver {
	return &CollisionSolver{
		cells:    make(map[[2]int][]ecs.Entity),
		contacts: make([]contactPairs, 0),
	}
}

func (cs *CollisionSolver) buildBroadPhase(ctx *core.Context) {
	ps := core.MustResource[*Settings](ctx)
	w := core.MustResource[*ecs.World](ctx)

	for i := range cs.cells {
		delete(cs.cells, i)
	}

	entities := w.Query(ecs.And(
		ecs.Has[*transform.Transform](),
		ecs.Has[*Aabb](),
	))

	for _, e := range entities {
		tr, _ := ecs.GetComponent[*transform.Transform](w, e)
		a, _ := ecs.GetComponent[*Aabb](w, e)

		minX := float64(tr.Position.X) - a.HalfWidth
		minY := float64(tr.Position.Y) - a.HalfHeight
		maxX := float64(tr.Position.X) + a.HalfWidth
		maxY := float64(tr.Position.Y) + a.HalfHeight

		cX0, cY0 := int(minX/ps.CellSize), int(minY/ps.CellSize)
		cX1, cY1 := int(maxX/ps.CellSize), int(maxY/ps.CellSize)

		for cY := cY0; cY <= cY1; cY++ {
			for cX := cX0; cX <= cX1; cX++ {
				key := [2]int{cX, cY}
				cs.cells[key] = append(cs.cells[key], e)
			}
		}

	}
}

func (cs *CollisionSolver) narrowPhaseAABB(ctx *core.Context) {
	w := core.MustResource[*ecs.World](ctx)

	cs.contacts = cs.contacts[:0]
	seen := core.NewHashSet[[2]ecs.Entity]()

	getBox := func(e ecs.Entity) (x, y, hw, hh float64, ok bool) {
		tr, ok1 := ecs.GetComponent[*transform.Transform](w, e)
		a, ok2 := ecs.GetComponent[*Aabb](w, e)

		if !ok1 || !ok2 {
			return
		}

		return float64(tr.Position.X), float64(tr.Position.Y), a.HalfWidth, a.HalfHeight, true
	}

	for _, list := range cs.cells {
		n := len(list)
		for i := range list {
			for j := i + 1; j < n; j++ {
				a, b := list[i], list[j]
				key := [2]ecs.Entity{a, b}

				if seen.Contains(key) {
					continue
				}

				seen.Add(key)

				aX, aY, aHw, aHh, ok1 := getBox(a)
				bX, bY, bHw, bHh, ok2 := getBox(b)

				if !ok1 || !ok2 {
					continue
				}
				dX := bX - aX
				pX := (aHw + bHw) - math.Abs(dX)
				if pX <= 0 {
					continue
				}

				dY := bY - aY
				pY := (aHh + bHh) - math.Abs(dY)
				if pY <= 0 {
					continue
				}

				cs.contacts = append(cs.contacts, contactPairs{
					a:    a,
					b:    b,
					penX: pX,
					penY: pY,
				})
			}
		}
	}

}

func (cs *CollisionSolver) resolveContactsAABB(ctx *core.Context) {
	w := core.MustResource[*ecs.World](ctx)

	for _, c := range cs.contacts {
		// components
		tra, _ := ecs.GetComponent[*transform.Transform](w, c.a)
		trb, _ := ecs.GetComponent[*transform.Transform](w, c.b)

		rba, _ := ecs.GetComponent[*RigidBody](w, c.a)
		rbb, _ := ecs.GetComponent[*RigidBody](w, c.b)

		va, _ := ecs.GetComponent[*Velocity](w, c.a)
		vb, _ := ecs.GetComponent[*Velocity](w, c.b)

		// TODO: collision filtering

		if c.penX < c.penY {
			dir := math.Copysign(1.0, float64(trb.Position.X-tra.Position.X))
			resolveX(tra, trb, rba, rbb, dir*c.penX*0.5, va, vb)
		} else {
			dir := math.Copysign(1.0, float64(trb.Position.Y-tra.Position.Y))
			resolveY(tra, trb, rba, rbb, dir*c.penY*0.5, va, vb)
		}
	}
}

func resolveX(tra, trb *transform.Transform, rba, rbb *RigidBody, mtv float64, va, vb *Velocity) {
	if rba.Type == Static && rbb.Type == Static {
		return
	}

	switch {
	case rba.Type != Static && rbb.Type != Static:
		tra.Position = tra.Position.Add(core.Vector[float32]{X: float32(-mtv), Y: 0})
		trb.Position = trb.Position.Add(core.Vector[float32]{X: float32(mtv), Y: 0})
		va.X = 0
		vb.X = 0
	case rbb.Type == Static:
		tra.Position = tra.Position.Add(core.Vector[float32]{X: float32(-mtv * 2), Y: 0})
		va.X = 0
	case rba.Type == Static:
		trb.Position = trb.Position.Add(core.Vector[float32]{X: float32(mtv * 2), Y: 0})
		vb.X = 0
	}

}

func resolveY(tra, trb *transform.Transform, rba, rbb *RigidBody, mtv float64, va, vb *Velocity) {
	if rba.Type == Static && rbb.Type == Static {
		return
	}

	switch {
	case rba.Type != Static && rbb.Type != Static:
		tra.Position = tra.Position.Add(core.Vector[float32]{Y: float32(-mtv), X: 0})
		trb.Position = trb.Position.Add(core.Vector[float32]{Y: float32(mtv), X: 0})
		va.Y = 0
		vb.Y = 0
	case rbb.Type == Static:
		tra.Position = tra.Position.Add(core.Vector[float32]{Y: float32(-mtv * 2), X: 0})
		va.Y = 0
	case rba.Type == Static:
		trb.Position = trb.Position.Add(core.Vector[float32]{Y: float32(mtv * 2), X: 0})
		vb.Y = 0
	}

}
