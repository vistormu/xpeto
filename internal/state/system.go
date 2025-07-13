package state

import (
	"time"

	"github.com/vistormu/xpeto/internal/ecs"
	"github.com/vistormu/xpeto/internal/event"
	st "github.com/vistormu/xpeto/internal/structures"
)

type System[S comparable] struct {
	requests *st.QueueArray[S]

	accumulator float64
	lastTime    time.Time
	fixedDelta  float64
}

func NewSystem[S comparable]() *System[S] {
	return &System[S]{
		requests:    st.NewQueueArray[S](),
		accumulator: 0,
		lastTime:    time.Now(),
		fixedDelta:  1.0 / 60.0,
	}
}

func (s *System[S]) OnEnter(ctx *ecs.Context) {
	em := ecs.MustResource[*event.Manager](ctx)

	event.Subscribe(em, func(data NextState[S]) {
		s.requests.Enqueue(data.State)
	})
}

func (s *System[S]) Update(ctx *ecs.Context) {
	// update time
	now := time.Now()
	frameTime := now.Sub(s.lastTime).Seconds()
	s.lastTime = now

	if frameTime > 0.25 {
		frameTime = 0.25
	}
	s.accumulator += frameTime

	steps := int(s.accumulator / s.fixedDelta)
	if steps > 0 {
		s.accumulator -= float64(steps) * s.fixedDelta
	}

	sm := ecs.MustResource[*Manager[S]](ctx)
	for !s.requests.IsEmpty() {
		state, _ := s.requests.Dequeue()

		// on exit functions
		onExitFns, ok := sm.onExitFns[state]
		if ok {
			for _, onExit := range onExitFns {
				onExit(ctx)
			}
		}

		// transition functions
		onTransitionFns, ok := sm.onTransitionFns[[2]S{sm.previous, state}]
		if ok {
			for _, onTransition := range onTransitionFns {
				onTransition(ctx)
			}
		}

		// on enter functions
		onEnterFns, ok := sm.onEnterFns[state]
		if ok {
			for _, onEnter := range onEnterFns {
				onEnter(ctx)
			}
		}
	}

	// update active states
	fixedUpdateFns, ok := sm.fixedUpdateFns[sm.active]
	if ok {
		for _, updateFn := range fixedUpdateFns {
			for i := 0; i < steps; i++ {
				updateFn(ctx, float32(s.fixedDelta))
			}
		}
	}

	updateFns, ok := sm.updateFns[sm.active]
	if ok {
		for _, updateFn := range updateFns {
			updateFn(ctx, float32(frameTime))
		}
	}
}
