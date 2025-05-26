package scene

import (
	"reflect"

	"github.com/vistormu/xpeto/ecs"
	"github.com/vistormu/xpeto/event"
	"github.com/vistormu/xpeto/image"
	st "github.com/vistormu/xpeto/internal/structures"
	"github.com/vistormu/xpeto/scene"
)

type action uint8

const (
	push action = iota
	pop
)

type transition struct {
	action action
	scene  reflect.Type
}

type System struct {
	transitions *st.QueueArray[transition]
	exiting     *st.QueueArray[scene.Scene]
	active      []scene.Scene
}

func NewSystem() *System {
	return &System{
		transitions: st.NewQueueArray[transition](),
		exiting:     st.NewQueueArray[scene.Scene](),
		active:      make([]scene.Scene, 0),
	}
}

func (s *System) OnLoad(ctx *ecs.Context) {
	em, _ := ecs.GetResource[*event.Manager](ctx)

	event.Subscribe[ScenePush](em, func(e ScenePush) {
		s.transitions.Enqueue(transition{
			action: push,
			scene:  e.Scene,
		})
	})

	event.Subscribe[ScenePop](em, func(e ScenePop) {
		s.transitions.Enqueue(transition{
			action: pop,
		})
	})
}

func (s *System) OnUnload(ctx *ecs.Context) {}

func (s *System) FixedUpdate(ctx *ecs.Context, dt float32) {}

func (s *System) Draw(screen *image.Image) {
	if len(s.active) == 0 {
		// log.Println("No active scenes to draw")
		return
	}

	// draw the first active scene TMP
	s.active[0].Draw(screen)
}

func (s *System) Update(ctx *ecs.Context, dt float32) {
	sm, _ := ecs.GetResource[*scene.Manager](ctx)

	// handle exiting scenes
	for !s.exiting.IsEmpty() {
		sc, _ := s.exiting.Dequeue()
		sc.OnUnload(ctx)
	}

	// handle transitions
	for !s.transitions.IsEmpty() {
		t, _ := s.transitions.Dequeue()

		switch t.action {
		case push:
			old, _ := sm.Current()
			new_, ok := sm.Scene(t.scene)
			if !ok || new_ == nil {
				continue
			}

			// unload old scene
			if old != nil {
				old.OnExit(ctx)
			}
			if !sm.IsActive(new_) {
				new_.OnLoad(ctx)
			}
			new_.OnEnter(ctx)

			sm.Push(t.scene)

		case pop:
			old, ok := sm.Current()
			if !ok || old == nil {
				continue
			}

			// unload old scene
			old.OnExit(ctx)
			s.exiting.Enqueue(old)
			sm.Pop()

			new_, ok := sm.Current()
			if !ok || new_ == nil {
				continue
			}

			// load new scene
			new_.OnEnter(ctx)
		}
	}

	// get active scenes
	s.active = sm.Active()
	if len(s.active) == 0 {
		// log.Println("No active scenes to update")
		return
	}

	// update the first active scene
	s.active[0].Update(ctx, dt)
}
