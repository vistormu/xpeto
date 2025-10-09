package state

import (
	"testing"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

type mockState int

const (
	firstState mockState = iota
	secondState
)

func TestStateMachine(t *testing.T) {
	// prepare
	w := ecs.NewWorld()
	sch := schedule.NewScheduler()

	// state machine
	AddStateMachine(sch, firstState)

	checkOrder := func(expected uint64) func(*ecs.World) {
		return func(w *ecs.World) {
			if expected != ecs.GetSystemId(w) {
				t.Fatalf("wrong order: expected %d, got %d", expected, ecs.GetSystemId(w))
			}
		}
	}

	stages := []schedule.StageFn{
		OnEnter(firstState),
		OnExit(firstState),
		OnTransition(firstState, secondState),
		OnEnter(secondState),
	}

	for i, s := range stages {
		schedule.AddSystem(sch, s, checkOrder(uint64(i+3))) // starts at 1 and the state machine registers two systems
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
