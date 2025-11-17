package state

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

// =====
// state
// =====
type state[T comparable] struct {
	current T
}

type nextState[T comparable] struct {
	next    *T
	pending bool
}

type transition[T comparable] struct {
	from, to T
}

// =============
// state machine
// =============
type stateMachine[T comparable] struct {
	initial      T
	onExit       map[T][]*schedule.Schedule
	onTransition map[transition[T]][]*schedule.Schedule
	onEnter      map[T][]*schedule.Schedule
}

func newStateMachine[T comparable](initial T) *stateMachine[T] {
	return &stateMachine[T]{
		initial:      initial,
		onExit:       make(map[T][]*schedule.Schedule),
		onTransition: make(map[transition[T]][]*schedule.Schedule),
		onEnter:      make(map[T][]*schedule.Schedule),
	}
}

func (sm *stateMachine[T]) addOnExit(st T, s *schedule.Schedule) {
	_, ok := sm.onExit[st]
	if !ok {
		sm.onExit[st] = make([]*schedule.Schedule, 0)
	}

	sm.onExit[st] = append(sm.onExit[st], s)
}

func (sm *stateMachine[T]) addOnTransition(from, to T, s *schedule.Schedule) {
	key := transition[T]{from: from, to: to}

	_, ok := sm.onTransition[key]
	if !ok {
		sm.onTransition[key] = make([]*schedule.Schedule, 0)
	}

	sm.onTransition[key] = append(sm.onTransition[key], s)
}

func (sm *stateMachine[T]) addOnEnter(st T, s *schedule.Schedule) {
	_, ok := sm.onEnter[st]
	if !ok {
		sm.onEnter[st] = make([]*schedule.Schedule, 0)
	}

	sm.onEnter[st] = append(sm.onEnter[st], s)
}

func (sm *stateMachine[T]) startup(w *ecs.World) {
	ecs.AddResource(w, state[T]{})
	ecs.AddResource(w, nextState[T]{next: &sm.initial, pending: true})
}

func (sm *stateMachine[T]) run(w *ecs.World, ss []*schedule.Schedule) {
	for _, sch := range ss {
		if sch == nil {
			continue
		}

		execute := true
		for _, c := range sch.Conditions {
			if c == nil {
				continue
			}

			execute = execute && c(w)
		}

		if execute {
			ecs.SetSystemInfo(w, sch.Id, "")
			sch.System(w)
		}
	}
}

func (sm *stateMachine[T]) update(w *ecs.World) {
	next, _ := ecs.GetResource[nextState[T]](w)

	if !next.pending || next.next == nil {
		return
	}

	current, _ := ecs.GetResource[state[T]](w)

	from := current.current
	to := *next.next
	next.pending = false
	next.next = nil

	if from == to {
		return
	}

	sm.run(w, sm.onExit[from])
	sm.run(w, sm.onTransition[transition[T]{from, to}])
	sm.run(w, sm.onEnter[to])

	current.current = to
}

// ===
// API
// ===
func AddStateMachine[T comparable](sch *schedule.Scheduler, initial T) {
	sm := newStateMachine(initial)
	schedule.AddExtra(sch, sm)
	schedule.AddSystem(sch, schedule.PreStartup, sm.startup)
	schedule.AddSystem(sch, schedule.StateTransition, sm.update)
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

	ns.next = &s
	ns.pending = true

	return true
}
