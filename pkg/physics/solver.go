package physics

import (
	"math"

	"github.com/vistormu/go-dsa/set"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/pkg/time"
	"github.com/vistormu/xpeto/core/pkg/transform"
)

func applyGravity(w *ecs.World) {
	ps, _ := ecs.GetResource[Settings](w)
	clk, _ := ecs.GetResource[time.FixedClock](w)

	q := ecs.NewQuery2[Velocity, RigidBody](w)

	for _, b := range q.Iter() {
		v := b.A()
		rb := b.B()

		if rb.Type != Dynamic {
			continue
		}

		scale := 1.0
		gs, ok := ecs.GetComponent[GravityScale](w, b.Entity())
		if ok {
			scale = gs.Value
		}

		v.Y += float64(ps.GravityY) * scale * clk.Timestep.Seconds()
		v.X += float64(ps.GravityX) * scale * clk.Timestep.Seconds()
	}
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

		tr.X += v.X * clk.Timestep.Seconds()
		tr.Y += v.Y * clk.Timestep.Seconds()
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

func (cs *CollisionSolver) buildBroadPhase(w *ecs.World) {
	ps, _ := ecs.GetResource[Settings](w)

	for i := range cs.cells {
		delete(cs.cells, i)
	}

	q := ecs.NewQuery2[Aabb, transform.Transform](w)

	for _, b := range q.Iter() {
		a := b.A()
		tr := b.B()

		minX := tr.X - a.HalfWidth
		minY := tr.Y - a.HalfHeight
		maxX := tr.X + a.HalfWidth
		maxY := tr.Y + a.HalfHeight

		cX0, cY0 := int(minX/ps.CellSize), int(minY/ps.CellSize)
		cX1, cY1 := int(maxX/ps.CellSize), int(maxY/ps.CellSize)

		for cY := cY0; cY <= cY1; cY++ {
			for cX := cX0; cX <= cX1; cX++ {
				key := [2]int{cX, cY}
				cs.cells[key] = append(cs.cells[key], b.Entity())
			}
		}

	}
}

func (cs *CollisionSolver) narrowPhaseAABB(w *ecs.World) {
	cs.contacts = cs.contacts[:0]
	seen := set.NewHashSet[[2]ecs.Entity]()

	getBox := func(e ecs.Entity) (x, y, hw, hh float64, ok bool) {
		tr, ok1 := ecs.GetComponent[transform.Transform](w, e)
		a, ok2 := ecs.GetComponent[Aabb](w, e)

		if !ok1 || !ok2 {
			return
		}

		return tr.X, tr.Y, a.HalfWidth, a.HalfHeight, true
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

func (cs *CollisionSolver) resolveContactsAABB(w *ecs.World) {
	for _, c := range cs.contacts {
		// components
		tra, _ := ecs.GetComponent[transform.Transform](w, c.a)
		trb, _ := ecs.GetComponent[transform.Transform](w, c.b)

		rba, _ := ecs.GetComponent[RigidBody](w, c.a)
		rbb, _ := ecs.GetComponent[RigidBody](w, c.b)

		va, _ := ecs.GetComponent[Velocity](w, c.a)
		vb, _ := ecs.GetComponent[Velocity](w, c.b)

		// TODO: collision filtering

		if c.penX < c.penY {
			dir := math.Copysign(1.0, trb.X-tra.X)
			resolveX(tra, trb, rba, rbb, dir*c.penX*0.5, va, vb)
		} else {
			dir := math.Copysign(1.0, trb.Y-tra.Y)
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
		tra.X -= mtv
		trb.X += mtv
		va.X = 0
		vb.X = 0
	case rbb.Type == Static:
		tra.X -= mtv * 2
		va.X = 0
	case rba.Type == Static:
		trb.X += mtv * 2
		vb.X = 0
	}

}

func resolveY(tra, trb *transform.Transform, rba, rbb *RigidBody, mtv float64, va, vb *Velocity) {
	if rba.Type == Static && rbb.Type == Static {
		return
	}

	switch {
	case rba.Type != Static && rbb.Type != Static:
		tra.Y -= mtv
		trb.Y += mtv
		va.Y = 0
		vb.Y = 0
	case rbb.Type == Static:
		tra.Y -= mtv * 2
		va.Y = 0
	case rba.Type == Static:
		trb.Y += mtv * 2
		vb.Y = 0
	}

}
