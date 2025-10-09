package schedule

import (
	"testing"

	"github.com/vistormu/xpeto/core/ecs"
)

func TestOrdering(t *testing.T) {
	w := ecs.NewWorld()
	sch := NewScheduler()

	checkOrder := func(expected uint64) func(*ecs.World) {
		return func(w *ecs.World) {
			if expected != ecs.GetSystemId(w) {
				t.Fatal("ids did not match")
			}
		}
	}

	stages := []StageFn{
		PreStartup,
		Startup,
		PostStartup,
		First,
		PreUpdate,
		FixedFirst,
		FixedPreUpdate,
		FixedUpdate,
		FixedPostUpdate,
		FixedLast,
		Update,
		PostUpdate,
		Last,
	}

	for i, s := range stages {
		AddSystem(sch, s, checkOrder(uint64(i+1)))
	}

	sch.RunStartup(w)
	sch.RunUpdate(w)
	sch.RunDraw(w)
}

type mockState int

const (
	firstState mockState = iota
	secondState
)

func TestStateMachine(t *testing.T) {
	// prepare
	w := ecs.NewWorld()
	sch := NewScheduler()

	// state machine
	AddStateMachine(sch, firstState)

	checkOrder := func(expected uint64) func(*ecs.World) {
		return func(w *ecs.World) {
			if expected != ecs.GetSystemId(w) {
				t.Fatalf("wrong order: expected %d, got %d", expected, ecs.GetSystemId(w))
			}
		}
	}

	stages := []StageFn{
		OnEnter(firstState),
		OnExit(firstState),
		OnTransition(firstState, secondState),
		OnEnter(secondState),
	}

	for i, s := range stages {
		AddSystem(sch, s, checkOrder(uint64(i+3))) // starts at 1 and the state machine registers two systems
	}

	sch.RunStartup(w)

	_, ok := GetState[mockState](w)
	if !ok {
		t.Fatal("state not found")
	}

	_, ok = ecs.GetResource[nextState[mockState]](w)
	if !ok {
		t.Fatal("next state not found")
	}

	// update
	sch.RunUpdate(w)

	current, ok := GetState[mockState](w)
	if current != firstState {
		t.Fatal("state machine did not transition to first state")
	}

	next, _ := ecs.GetResource[nextState[mockState]](w)
	if next.next != nil || next.pending {
		t.Fatal("next state did not clean up")
	}

	// transition
	ok = SetNextState(w, secondState)
	if !ok {
		t.Fatal("could not set next state")
	}

	sch.RunUpdate(w)
}

func TestOnceAndOnceWhen(t *testing.T) {
	w := ecs.NewWorld()
	sch := NewScheduler()
	ran := 0

	AddSystem(sch, Update, func(*ecs.World) { ran++ }).RunIf(Once())
	sch.RunUpdate(w)
	sch.RunUpdate(w)
	if ran != 1 {
		t.Fatalf("Once ran %d times", ran)
	}

	ran = 0
	gate := false
	AddSystem(sch, Update, func(*ecs.World) { ran++ }).
		RunIf(OnceWhen(func(*ecs.World) bool { return gate }))

	sch.RunUpdate(w) // gate=false → not yet
	gate = true
	sch.RunUpdate(w) // fires now
	sch.RunUpdate(w) // no more
	if ran != 1 {
		t.Fatalf("OnceWhen ran %d times", ran)
	}
}
