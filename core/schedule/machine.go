package schedule

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"
	"github.com/vistormu/xpeto/internal/event"
)

type StateMachine[T comparable] struct {
	initial      T
	onExit       map[T][]ecs.System
	onTransition map[transition[T]][]ecs.System
	onEnter      map[T][]ecs.System
}

func newStateMachine[T comparable](initial T) *StateMachine[T] {
	return &StateMachine[T]{
		initial:      initial,
		onExit:       make(map[T][]ecs.System),
		onTransition: make(map[transition[T]][]ecs.System),
		onEnter:      make(map[T][]ecs.System),
	}
}

func (sm *StateMachine[T]) startup(ctx *core.Context) {
	core.AddResource(ctx, &State[T]{})
	core.AddResource(ctx, &NextState[T]{next: &sm.initial, pending: true})
}

func (sm *StateMachine[T]) update(ctx *core.Context) {
	next := core.MustResource[*NextState[T]](ctx)

	if !next.pending || next.next == nil {
		return
	}

	current := core.MustResource[*State[T]](ctx)

	from := current.current
	to := *next.next
	next.pending = false
	next.next = nil

	if from == to {
		return
	}

	// publish event
	eb := core.MustResource[*event.Bus](ctx)
	event.Publish(eb, EventStateTransition[T]{Exited: &from, Entered: &to})

	// on exit
	for _, sys := range sm.onExit[from] {
		sys(ctx)
	}

	// on transition
	for _, sys := range sm.onTransition[transition[T]{from: from, to: to}] {
		sys(ctx)
	}

	// on enter
	for _, sys := range sm.onEnter[to] {
		sys(ctx)
	}

	current.current = to
}
