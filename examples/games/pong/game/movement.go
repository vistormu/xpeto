package game

import (
	"math/rand/v2"

	"github.com/vistormu/xpeto"
	"github.com/vistormu/xpeto/pkg/physics"
)

func movePaddles(w *xp.World) {
	ev, ok := xp.GetEvents[MoveIntent](w)
	if !ok {
		return
	}

	c, _ := xp.GetResource[Config](w)
	q := xp.Query2[Paddle, physics.Velocity](w)
	for _, b := range q.Iter() {
		p, v := b.Components()

		v.Y = 0

		if len(ev) == 0 {
			continue
		}

		for _, e := range ev {
			if p.IsLeft && e.IsLeft && e.IsUp {
				v.Y -= c.PaddleSpeed
			} else if p.IsLeft && e.IsLeft && !e.IsUp {
				v.Y += c.PaddleSpeed
			}

			if !p.IsLeft && !e.IsLeft && e.IsUp {
				v.Y -= c.PaddleSpeed
			} else if !p.IsLeft && !e.IsLeft && !e.IsUp {
				v.Y += c.PaddleSpeed
			}
		}
	}
}

func moveBall(w *xp.World) {
	c, _ := xp.GetResource[Config](w)

	b, _ := xp.Query1[physics.Velocity](w, xp.With[Ball]()).First()
	v := b.Components()
	v.Y = rand.NormFloat64()*0.5 + 50
	v.X = rand.NormFloat64()*1 + (c.MaxBallSpeed+c.MinBallSpeed)/2
}

func ballTrail(w *xp.World) {
	b, _ := xp.Query2[xp.Transform, xp.PathShape](w, xp.With[Ball]()).Single()
	r, p := b.Components()
	p.AddPoint(float32(r.X), float32(r.Y))
}

func movementMiniPkg(_ *xp.World, sch *xp.Scheduler) {
	xp.AddSystem(sch, xp.OnEnter(statePlaying), moveBall)
	xp.AddSystem(sch, xp.Update, ballTrail).RunIf(xp.InState(statePlaying))
	xp.AddSystem(sch, xp.FixedUpdate, movePaddles).RunIf(xp.InState(statePlaying))
}
