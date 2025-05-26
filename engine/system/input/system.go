package input

import (
	"github.com/vistormu/xpeto/ecs"
	"github.com/vistormu/xpeto/event"
	"github.com/vistormu/xpeto/image"
	"github.com/vistormu/xpeto/input"
	st "github.com/vistormu/xpeto/internal/structures"
)

type System struct {
	pressedKeys *st.HashSet[input.Key]
}

func NewSystem() ecs.System {
	return &System{
		pressedKeys: st.NewHashSet[input.Key](),
	}
}

func (s *System) OnLoad(ctx *ecs.Context)                  {}
func (s *System) OnUnload(ctx *ecs.Context)                {}
func (s *System) FixedUpdate(ctx *ecs.Context, dt float32) {}
func (s *System) Draw(*image.Image)                        {}

func (is *System) Update(ctx *ecs.Context, dt float32) {
	im, _ := ecs.GetResource[*input.Manager](ctx)
	em, _ := ecs.GetResource[*event.Manager](ctx)

	// keys
	for _, key := range im.Keys() {
		// key presses
		if im.IsPressed(key) && !is.pressedKeys.Contains(key) {
			is.pressedKeys.Add(key)
			em.Publish(KeyPress{Key: key})
		}

		// key maintains
		if is.pressedKeys.Contains(key) && im.IsPressed(key) {
			em.Publish(KeyMaintain{Key: key})
		}

		// key releases
		if is.pressedKeys.Contains(key) && !im.IsPressed(key) {
			is.pressedKeys.Remove(key)
			em.Publish(KeyRelease{Key: key})
		}
	}
}
