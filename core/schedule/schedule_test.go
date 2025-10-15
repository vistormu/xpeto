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

func TestOrderingWithinStage(t *testing.T) {
	w := ecs.NewWorld()
	sch := NewScheduler()

	got := make([]string, 0, 3)
	sys := func(name string) ecs.System {
		return func(*ecs.World) { got = append(got, name) }
	}

	AddSystem(sch, Update, sys("A")).Label("A")
	AddSystem(sch, Update, sys("B")).Label("B")
	AddSystem(sch, Update, sys("C")).Label("C")

	sch.After("A")
	sch.Before("C")

	sch = NewScheduler()
	AddSystem(sch, Update, sys("A")).Label("A")
	AddSystem(sch, Update, sys("B")).Label("B").Before("C")
	AddSystem(sch, Update, sys("C")).Label("C").After("A")

	sch.RunUpdate(w)

	want := []string{"A", "B", "C"}
	if len(got) != len(want) {
		t.Fatalf("len mismatch: got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("order mismatch at %d: got %q want %q", i, got[i], want[i])
		}
	}
}
