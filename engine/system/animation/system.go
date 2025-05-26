package animation

import (
	"github.com/vistormu/xpeto/ecs"
	"github.com/vistormu/xpeto/engine/component"
	"github.com/vistormu/xpeto/image"
)

type System struct {
}

func NewSystem() ecs.System {
	return &System{}
}

func (s *System) OnLoad(*ecs.Context)               {}
func (s *System) OnUnload(*ecs.Context)             {}
func (s *System) FixedUpdate(*ecs.Context, float32) {}
func (s *System) Draw(*image.Image)                 {}

func (s *System) Update(ctx *ecs.Context, dt float32) {
	em, _ := ecs.GetResource[*ecs.EntityManager](ctx)

	entities := em.Query(ecs.And(
		ecs.Has[*component.Animation](),
		ecs.Has[*component.Renderable](),
	))

	for _, e := range entities {
		anim, _ := ecs.GetComponent[*component.Animation](em, e)
		renderable, _ := ecs.GetComponent[*component.Renderable](em, e)

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
		renderable.Sprite = anim.Frames[anim.Current]
	}
}
