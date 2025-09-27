package animation

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"
)

func Update(ctx *core.Context, dt float32) {
	w, _ := core.GetResource[*ecs.World](ctx)

	entities := w.Query(ecs.And(
		ecs.Has[*Animation](),
		ecs.Has[*render.Renderable](),
	))

	for _, e := range entities {
		anim, _ := ecs.GetComponent[*Animation](w, e)
		renderable, _ := ecs.GetComponent[*render.Renderable](w, e)

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
