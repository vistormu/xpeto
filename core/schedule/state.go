package schedule

import (
	"github.com/vistormu/go-dsa/hashmap"
	"github.com/vistormu/xpeto/core/ecs"
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
	onExit       map[T][]*schedule
	onTransition map[transition[T]][]*schedule
	onEnter      map[T][]*schedule
}

func newStateMachine[T comparable](initial T) *stateMachine[T] {
	return &stateMachine[T]{
		initial:      initial,
		onExit:       make(map[T][]*schedule),
		onTransition: make(map[transition[T]][]*schedule),
		onEnter:      make(map[T][]*schedule),
	}
}

func (sm *stateMachine[T]) addOnExit(st T, s *schedule) {
	_, ok := sm.onExit[st]
	if !ok {
		sm.onExit[st] = make([]*schedule, 0)
	}

	sm.onExit[st] = append(sm.onExit[st], s)
}

func (sm *stateMachine[T]) addOnTransition(from, to T, s *schedule) {
	key := transition[T]{from: from, to: to}

	_, ok := sm.onTransition[key]
	if !ok {
		sm.onTransition[key] = make([]*schedule, 0)
	}

	sm.onTransition[key] = append(sm.onTransition[key], s)
}

func (sm *stateMachine[T]) addOnEnter(st T, s *schedule) {
	_, ok := sm.onEnter[st]
	if !ok {
		sm.onEnter[st] = make([]*schedule, 0)
	}

	sm.onEnter[st] = append(sm.onEnter[st], s)
}

func (sm *stateMachine[T]) startup(w *ecs.World) {
	ecs.AddResource(w, state[T]{})
	ecs.AddResource(w, nextState[T]{next: &sm.initial, pending: true})
}

func (sm *stateMachine[T]) run(w *ecs.World, ss []*schedule) {
	for _, sch := range ss {
		if sch == nil {
			continue
		}

		execute := true
		for _, c := range sch.conditions {
			if c == nil {
				continue
			}

			execute = execute && c(w)
		}

		if execute {
			ecs.SetSystemId(w, sch.id)
			sch.system(w)
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

func AddStateMachine[T comparable](sch *Scheduler, initial T) {
	sm := newStateMachine(initial)
	hashmap.Add(sch.stateMachines, sm)

	// startup schedule
	s := newSchedule()
	s.stage = postStartup
	s.system = sm.startup
	s.id = sch.nextId
	sch.nextId++

	sch.addSchedule(s)

	// update schedule
	s = newSchedule()
	s.stage = stateTransition
	s.system = sm.update
	s.id = sch.nextId
	sch.nextId++

	sch.addSchedule(s)
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
