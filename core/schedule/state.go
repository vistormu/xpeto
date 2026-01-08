package schedule

import (
	"github.com/vistormu/xpeto/core/ecs"
)

// =====
// state
// =====
type state[T comparable] struct {
	current T
}

type nextState[T comparable] struct {
	next    T
	pending bool
}

// ==========
// transition
// ==========
type transition[T comparable] struct {
	from, to T
}

type transitionEvent struct {
	onExit       []uint64
	onTransition []uint64
	onEnter      []uint64
}

func newTransitionEvent() transitionEvent {
	return transitionEvent{
		onExit:       make([]uint64, 0),
		onTransition: make([]uint64, 0),
		onEnter:      make([]uint64, 0),
	}
}

func (te *transitionEvent) clear() {
	te.onExit = te.onExit[:0]
	te.onTransition = te.onTransition[:0]
	te.onEnter = te.onEnter[:0]
}

// =============
// state machine
// =============
type stateMachine[T comparable] struct {
	initial      T
	onExit       map[T][]uint64
	onTransition map[transition[T]][]uint64
	onEnter      map[T][]uint64
}

func newStateMachine[T comparable](initial T) stateMachine[T] {
	return stateMachine[T]{
		initial:      initial,
		onExit:       make(map[T][]uint64),
		onTransition: make(map[transition[T]][]uint64),
		onEnter:      make(map[T][]uint64),
	}
}

func (sm *stateMachine[T]) add(from, to *T, id uint64) {
	if from != nil && to == nil {
		sm.onExit[*from] = append(sm.onExit[*from], id)
	}

	if from != nil && to != nil {
		key := transition[T]{from: *from, to: *to}
		sm.onTransition[key] = append(sm.onTransition[key], id)
	}

	if from == nil && to != nil {
		sm.onEnter[*to] = append(sm.onEnter[*to], id)
	}
}

// =======
// systems
// =======
func (sm *stateMachine[T]) startup(w *ecs.World) {
	ecs.AddResource(w, state[T]{sm.initial})
	ecs.AddResource(w, nextState[T]{})
}

func (sm *stateMachine[T]) update(w *ecs.World) {
	next, ok := ecs.GetResource[nextState[T]](w)
	if !ok {
		ecs.AddResource(w, nextState[T]{})
		return
	}
	if !next.pending {
		return
	}

	current, ok := ecs.GetResource[state[T]](w)
	if !ok {
		ecs.AddResource(w, state[T]{sm.initial})
		return
	}

	from := current.current
	to := next.next

	next.pending = false

	if from == to {
		return
	}

	current.current = to

	tr, ok := ecs.GetResource[transitionEvent](w)
	if !ok {
		return
	}

	tr.onExit = append(tr.onExit, sm.onExit[from]...)
	tr.onTransition = append(tr.onTransition, sm.onTransition[transition[T]{from, to}]...)
	tr.onEnter = append(tr.onEnter, sm.onEnter[to]...)
}

// ===
// API
// ===
func AddStateMachine[T comparable](sch *Scheduler, initial T) {
	sm := newStateMachine(initial)
	addStateMachine(sch.store, sm)

	AddSystem(sch, PreStartup, sm.startup)
	AddSystem(sch, stateTransitionStage, sm.update)
}

func GetState[T comparable](w *ecs.World) (T, bool) {
	s, ok := ecs.GetResource[state[T]](w)
	return s.current, ok
}

func SetNextState[T comparable](w *ecs.World, s T) bool {
	ns, ok := ecs.GetResource[nextState[T]](w)
	if !ok {
		return false
	}

	ns.next = s
	ns.pending = true

	return true
}
