package animation

import (
	"github.com/vistormu/xpeto/internal/ecs"
	"github.com/vistormu/xpeto/pkg/render"
)

type System struct {
}

func NewSystem() *System {
	return &System{}
}

func (s *System) Update(ctx *ecs.Context, dt float32) {
	em, _ := ecs.GetResource[*ecs.Manager](ctx)

	entities := em.Query(ecs.And(
		ecs.Has[*Animation](),
		ecs.Has[*render.Renderable](),
	))

	for _, e := range entities {
		anim, _ := ecs.GetComponent[*Animation](em, e)
		renderable, _ := ecs.GetComponent[*render.Renderable](em, e)

		// Advance elapsed time
		anim.Elapsed += dt
		if anim.Elapsed >= anim.Duration {
			anim.Elapsed -= anim.Duration
			anim.Current++
			// Loop or clamp
			if anim.Current >= uint64(len(anim.Frames)) {
				if anim.Loop {
					anim.Current = 0
				} else {
					anim.Current = uint64(len(anim.Frames)) - 1
				}
			}
		}

		// Fetch sprite and update its image
		renderable.Image = anim.Frames[anim.Current]
	}
}
