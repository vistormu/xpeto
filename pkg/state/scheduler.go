package state

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/event"
)

type Scheduler[S comparable] struct {
	requests *core.QueueArray[S]

	accumulator float64
	lastTime    time.Time
	fixedDelta  float64
}

func NewScheduler[S comparable](initial S) *Scheduler[S] {
	return &Scheduler[S]{
		requests:    core.NewQueueArray[S](),
		accumulator: 0,
		lastTime:    time.Now(),
		fixedDelta:  1.0 / 60.0,
	}
}

func (s *Scheduler[S]) OnEnter(ctx *core.Context) {
	em := core.MustResource[*event.Bus](ctx)

	event.Subscribe(em, func(data NextState[S]) {
		s.requests.Enqueue(data.State)
	})
}

func (s *Scheduler[S]) Update(ctx *core.Context) {
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

	sm := core.MustResource[*Fsm[S]](ctx)
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
			for range steps {
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

func (s *Scheduler[S]) Draw(screen *ebiten.Image) {}
