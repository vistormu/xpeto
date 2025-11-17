package physics

import (
	"math"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/pkg/transform"
)

// =========
// rect-rect
// =========
type rectOBB struct {
	cx, cy float64
	ux, uy float64
	vx, vy float64
	ex, ey float64
}

func obbFromRect(r Rect, tr *transform.Transform) rectOBB {
	c := math.Cos(tr.Rotation)
	s := math.Sin(tr.Rotation)

	return rectOBB{
		cx: tr.X, cy: tr.Y,
		ux: c, uy: s,
		vx: -s, vy: c,
		ex: r.HalfW, ey: r.HalfH,
	}
}

func projectOBB(o rectOBB, ax, ay float64) (min, max float64) {
	pc := o.cx*ax + o.cy*ay
	r := o.ex*math.Abs(o.ux*ax+o.uy*ay) + o.ey*math.Abs(o.vx*ax+o.vy*ay)

	return pc - r, pc + r
}

func intervalOverlap(a0, a1, b0, b1 float64) (float64, bool) {
	if a0 > b1 || b0 > a1 {
		return 0, false
	}

	min1 := a1
	if b1 < min1 {
		min1 = b1
	}

	max0 := a0
	if b0 > max0 {
		max0 = b0
	}

	return min1 - max0, true
}

func satRectRect(a Rect, ta *transform.Transform, b Rect, tb *transform.Transform) (bool, float64, float64, float64) {
	oa := obbFromRect(a, ta)
	ob := obbFromRect(b, tb)

	axes := [][2]float64{
		{oa.ux, oa.uy}, {oa.vx, oa.vy},
		{ob.ux, ob.uy}, {ob.vx, ob.vy},
	}

	dx, dy := ob.cx-oa.cx, ob.cy-oa.cy

	minOverlap := math.Inf(1)
	nx, ny := 0.0, 0.0

	for _, ax := range axes {
		axX, axY := ax[0], ax[1]
		a0, a1 := projectOBB(oa, axX, axY)
		b0, b1 := projectOBB(ob, axX, axY)

		ov, ok := intervalOverlap(a0, a1, b0, b1)
		if !ok {
			return false, 0, 0, 0
		}

		if ov < minOverlap {
			minOverlap = ov
			if axX*dx+axY*dy < 0 {
				nx, ny = -axX, -axY
			} else {
				nx, ny = axX, axY
			}
		}
	}

	return true, nx, ny, minOverlap
}

func narrowPhaseRectRect(w *ecs.World) {
	s, _ := ecs.GetResource[Space](w)
	if s == nil || len(s.candidates) == 0 {
		return
	}

	for _, p := range s.candidates {
		ca, okA := ecs.GetComponent[Collider](w, p.A)
		cb, okB := ecs.GetComponent[Collider](w, p.B)
		if !okA || !okB {
			continue
		}

		// only rect-rect for now
		ra, okRA := ca.Shape.(Rect)
		rb, okRB := cb.Shape.(Rect)
		if !okRA || !okRB {
			continue
		}

		ta, okTA := ecs.GetComponent[transform.Transform](w, p.A)
		tb, okTB := ecs.GetComponent[transform.Transform](w, p.B)
		if !okTA || !okTB {
			continue
		}

		hit, nx, ny, depth := satRectRect(ra, ta, rb, tb)
		if !hit || depth <= 0 {
			continue
		}

		// simple penetration vector from A into free space
		penX, penY := nx*depth, ny*depth

		s.Contacts = append(s.Contacts, ContactPair{
			A: p.A, B: p.B,
			PenX: penX, PenY: penY,
			NormalX: nx, NormalY: ny,
			Depth: depth,
		})
	}
}

// ========
// contacts
// ========
const (
	penetrationSlop   = 0.01 // tolerate tiny overlaps
	correctionPercent = 0.8  // project 80% per step for stability
	solverIterations  = 1    // bump to 4+ for stacking later
)

func invMassOf(rb *RigidBody) float64 {
	if rb == nil {
		// No component ⇒ dynamic with unit mass
		return 1.0
	}
	switch rb.Type {
	case Static, Kinematic:
		return 0.0
	case Dynamic:
		m := rb.Mass
		if !(m > 0) { // NaN or ≤0
			m = 1.0
		}
		return 1.0 / m
	default:
		return 1.0
	}
}

func resolveContacts(w *ecs.World) {
	s, _ := ecs.GetResource[Space](w)
	if s == nil || len(s.Contacts) == 0 {
		return
	}

	for it := 0; it < solverIterations; it++ {
		for _, c := range s.Contacts {
			ca, okA := ecs.GetComponent[Collider](w, c.A)
			cb, okB := ecs.GetComponent[Collider](w, c.B)
			if !okA || !okB {
				continue
			}
			// sensors: report but don't separate
			if ca.Sensor || cb.Sensor {
				continue
			}

			ta, okTA := ecs.GetComponent[transform.Transform](w, c.A)
			tb, okTB := ecs.GetComponent[transform.Transform](w, c.B)
			if !okTA || !okTB {
				continue
			}

			rba, _ := ecs.GetComponent[RigidBody](w, c.A)
			rbb, _ := ecs.GetComponent[RigidBody](w, c.B)

			invA := invMassOf(rba)
			invB := invMassOf(rbb)
			if invA == 0 && invB == 0 {
				continue
			} // both immovable

			// MTV with slop and percentage
			over := c.Depth - penetrationSlop
			if over <= 0 {
				continue
			}

			corrMag := over * correctionPercent
			corrX := c.NormalX * corrMag
			corrY := c.NormalY * corrMag

			// distribute by inverse mass
			invSum := invA + invB
			wa, wb := 0.0, 0.0
			if invSum > 0 {
				wa = invA / invSum
				wb = invB / invSum
			}

			// Move A opposite the normal, B along the normal.
			// Static/Kinematic (inv=0) stay put; Dynamics move.
			if invA > 0 {
				ta.X -= corrX * wa
				ta.Y -= corrY * wa
			}
			if invB > 0 {
				tb.X += corrX * wb
				tb.Y += corrY * wb
			}
		}
	}
}

func combineAverage(a, b float64) float64 {
	return 0.5 * (a + b)
}

func resolveContactsImpulses(w *ecs.World) {
	s, _ := ecs.GetResource[Space](w)
	if s == nil || len(s.Contacts) == 0 {
		return
	}

	for _, c := range s.Contacts {
		// Colliders / bodies
		ca, okA := ecs.GetComponent[Collider](w, c.A)
		cb, okB := ecs.GetComponent[Collider](w, c.B)
		if !okA || !okB {
			continue
		}
		// Sensors: skip impulses
		if ca.Sensor || cb.Sensor {
			continue
		}

		rba, _ := ecs.GetComponent[RigidBody](w, c.A)
		rbb, _ := ecs.GetComponent[RigidBody](w, c.B)
		invA := invMassOf(rba)
		invB := invMassOf(rbb)
		if invA == 0 && invB == 0 {
			continue
		} // both immovable

		// Transforms (we don't actually need positions here, but we keep the fetch
		// pattern consistent; later for angular we'll need contact points)
		_, okTA := ecs.GetComponent[transform.Transform](w, c.A)
		_, okTB := ecs.GetComponent[transform.Transform](w, c.B)
		if !okTA || !okTB {
			continue
		}

		// Velocities (default to zero if missing)
		va, _ := ecs.GetComponent[Velocity](w, c.A)
		vb, _ := ecs.GetComponent[Velocity](w, c.B)
		var vax, vay, vbx, vby float64
		if va != nil {
			vax, vay = va.X, va.Y
		}
		if vb != nil {
			vbx, vby = vb.X, vb.Y
		}

		// Relative velocity
		rvx, rvy := vbx-vax, vby-vay

		// Contact normal (assumed unit-length from narrow phase)
		nx, ny := c.NormalX, c.NormalY

		// Normal relative speed
		vn := rvx*nx + rvy*ny
		if vn > 0 {
			// Already separating along the normal; no bounce / no friction needed.
			continue
		}

		// Effective mass (no rotation yet)
		keff := invA + invB
		if keff == 0 {
			continue
		}

		// Combine coefficients (Average like Rapier's default)
		e := combineAverage(safeRestitution(rba), safeRestitution(rbb))
		mu := combineAverage(safeFriction(rba), safeFriction(rbb))

		// --- Normal impulse (restitution) ---
		jn := -(1.0 + e) * vn / keff
		impNx, impNy := jn*nx, jn*ny

		// Apply normal impulse
		if invA > 0 {
			vax -= impNx * invA
			vay -= impNy * invA
		}
		if invB > 0 {
			vbx += impNx * invB
			vby += impNy * invB
		}

		// Recompute relative velocity after normal impulse (for friction)
		rvx, rvy = vbx-vax, vby-vay

		// --- Tangent (Coulomb friction) ---
		// Tangent = reject(v, n) normalized
		vtX, vtY := rvx-(rvx*nx+rvy*ny)*nx, rvy-(rvx*nx+rvy*ny)*ny
		vtLen := math.Hypot(vtX, vtY)

		if vtLen > 1e-8 && jn > 0 {
			tx, ty := vtX/vtLen, vtY/vtLen
			vt := rvx*tx + rvy*ty // tangential relative speed

			jt := -vt / keff
			// Coulomb cone clamp
			maxFric := mu * jn
			if jt > maxFric {
				jt = maxFric
			}
			if jt < -maxFric {
				jt = -maxFric
			}

			impTx, impTy := jt*tx, jt*ty

			// Apply friction impulse
			if invA > 0 {
				vax -= impTx * invA
				vay -= impTy * invA
			}
			if invB > 0 {
				vbx += impTx * invB
				vby += impTy * invB
			}
		}

		// Write back velocities (create component if absent)
		if va == nil {
			ecs.AddComponent(w, c.A, Velocity{X: vax, Y: vay})
		} else {
			va.X, va.Y = vax, vay
		}
		if vb == nil {
			ecs.AddComponent(w, c.B, Velocity{X: vbx, Y: vby})
		} else {
			vb.X, vb.Y = vbx, vby
		}
	}
}

func safeRestitution(rb *RigidBody) float64 {
	if rb == nil {
		return 0
	}
	return rb.Restitution
}

func safeFriction(rb *RigidBody) float64 {
	if rb == nil {
		return 0.5
	}
	return rb.Friction
}
