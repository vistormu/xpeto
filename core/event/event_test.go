package event

import (
	"testing"

	"github.com/vistormu/xpeto/core/ecs"
)

type evA struct{ v int }
type evB struct{ s string }

func TestFirstReadReturnsAllSoFar(t *testing.T) {
	w := ecs.NewWorld()
	ecs.AddResource(w, newBus())

	ecs.SetSystemId(w, 1)

	AddEvent(w, evA{1})
	AddEvent(w, evA{2})

	got, ok := GetEvents[evA](w)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if len(got) != 2 || got[0].v != 1 || got[1].v != 2 {
		t.Fatalf("unexpected first read: %#v", got)
	}

	got, ok = GetEvents[evA](w)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if len(got) != 0 {
		t.Fatalf("expected empty second read; got %#v", got)
	}
}

func TestDoubleBufferingAcrossUpdate(t *testing.T) {
	w := ecs.NewWorld()
	ecs.AddResource(w, newBus())

	ecs.SetSystemId(w, 1)

	AddEvent(w, evA{10})
	Update(w)

	AddEvent(w, evA{20})
	AddEvent(w, evA{30})

	got, ok := GetEvents[evA](w)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if len(got) != 3 || got[0].v != 10 || got[1].v != 20 || got[2].v != 30 {
		t.Fatalf("unexpected read across update: %#v", got)
	}
	got, _ = GetEvents[evA](w)
	if len(got) != 0 {
		t.Fatalf("expected empty after consuming; got %#v", got)
	}
}

func TestPerSystemIsolation(t *testing.T) {
	w := ecs.NewWorld()
	ecs.AddResource(w, newBus())

	ecs.SetSystemId(w, 1)
	AddEvent(w, evA{1})
	got1, _ := GetEvents[evA](w)
	if len(got1) != 1 || got1[0].v != 1 {
		t.Fatalf("sys1 wrong first read: %#v", got1)
	}

	ecs.SetSystemId(w, 2)
	got2, _ := GetEvents[evA](w)
	if len(got2) != 1 || got2[0].v != 1 {
		t.Fatalf("sys2 should independently see event: %#v", got2)
	}

	ecs.SetSystemId(w, 1)
	AddEvent(w, evA{2})
	got1, _ = GetEvents[evA](w)
	if len(got1) != 1 || got1[0].v != 2 {
		t.Fatalf("sys1 should see only new event: %#v", got1)
	}

	ecs.SetSystemId(w, 2)
	got2, _ = GetEvents[evA](w)
	if len(got2) != 1 || got2[0].v != 2 {
		t.Fatalf("sys2 should now see the second event: %#v", got2)
	}
}

func TestPerTypeIsolationSameSystem(t *testing.T) {
	w := ecs.NewWorld()
	ecs.AddResource(w, newBus())
	ecs.SetSystemId(w, 1)

	AddEvent(w, evA{7})
	AddEvent(w, evB{"x"})

	as, ok := GetEvents[evA](w)
	if !ok || len(as) != 1 || as[0].v != 7 {
		t.Fatalf("unexpected A read: %#v", as)
	}

	bs, ok := GetEvents[evB](w)
	if !ok || len(bs) != 1 || bs[0].s != "x" {
		t.Fatalf("unexpected B read: %#v", bs)
	}
}
