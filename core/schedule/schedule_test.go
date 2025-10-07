package schedule

import (
	"testing"

	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/event"
)

func TestOrdering(t *testing.T) {
	ctx := core.NewContext()
	sch := NewScheduler()

	AddSystem(sch, PreStartup, func(*core.Context) {
		t.Log("0")
	})
	AddSystem(sch, Startup, func(*core.Context) {
		t.Log("1")
	})
	AddSystem(sch, PostStartup, func(*core.Context) {
		t.Log("2")
	})

	AddSystem(sch, First, func(*core.Context) {
		t.Log("3")
	})
	AddSystem(sch, PreUpdate, func(*core.Context) {
		t.Log("4")
	})

	AddSystem(sch, FixedFirst, func(*core.Context) {
		t.Log("5")
	})
	AddSystem(sch, FixedPreUpdate, func(*core.Context) {
		t.Log("6")
	})
	AddSystem(sch, FixedUpdate, func(*core.Context) {
		t.Log("7")
	})
	AddSystem(sch, FixedPostUpdate, func(*core.Context) {
		t.Log("8")
	})
	AddSystem(sch, FixedLast, func(*core.Context) {
		t.Log("9")
	})

	AddSystem(sch, Update, func(*core.Context) {
		t.Log("10")
	})
	AddSystem(sch, PostUpdate, func(*core.Context) {
		t.Log("11")
	})
	AddSystem(sch, Last, func(*core.Context) {
		t.Log("12")
	})

	AddSystem(sch, PreDraw, func(*core.Context) {
		t.Log("13")
	})
	AddSystem(sch, Draw, func(*core.Context) {
		t.Log("14")
	})
	AddSystem(sch, PostDraw, func(*core.Context) {
		t.Log("15")
	})

	sch.RunStartup(ctx)
	sch.RunUpdate(ctx)
	sch.RunDraw(ctx)
}

type states int

const (
	firstState states = iota + 1
	secondState
)

func TestStateMachine(t *testing.T) {
	// prepare
	ctx := core.NewContext()
	eb := event.NewBus()
	core.AddResource(ctx, eb)

	sch := NewScheduler()

	// state machine
	AddStateMachine(sch, firstState)

	AddSystem(sch, OnEnter(firstState), func(*core.Context) {
		t.Log("first state entered")
	})
	AddSystem(sch, OnExit(firstState), func(*core.Context) {
		t.Log("first state exited")
	})
	AddSystem(sch, OnTransition(firstState, secondState), func(*core.Context) {
		t.Log("transitioned from first to second state")
	})
	AddSystem(sch, OnEnter(secondState), func(*core.Context) {
		t.Log("second state entered")
	})

	// startup
	sch.RunStartup(ctx)

	_, ok := core.GetResource[*State[states]](ctx)
	if !ok {
		t.Fatal("state not found")
	}

	_, ok = core.GetResource[*NextState[states]](ctx)
	if !ok {
		t.Fatal("next state not found")
	}

	// update
	sch.RunUpdate(ctx)

	current := core.MustResource[*State[states]](ctx)
	if current.Get() != firstState {
		t.Fatal("state machine did not transition to first state")
	}

	next := core.MustResource[*NextState[states]](ctx)
	if next.next != nil || next.pending {
		t.Fatal("next state did not clean up")
	}

	// transition
	next.Set(secondState)
	sch.RunUpdate(ctx)
}
