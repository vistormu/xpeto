package event

import (
	"sync"
	"testing"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

type evA struct{ v int }
type evB struct{ s string }

func newTestWorld() *ecs.World {
	w := ecs.NewWorld()
	ecs.AddResource(w, newBus())
	ecs.AddResource(w, schedule.RunningSystem{})
	return w
}

func setSystemInfo(w *ecs.World, id uint64, label string) {
	rs, _ := ecs.GetResource[schedule.RunningSystem](w)
	rs.Id = id
	rs.Label = label
}

func TestFirstReadReturnsAllSoFar(t *testing.T) {
	w := newTestWorld()
	setSystemInfo(w, 1, "sys")

	AddEvent(w, evA{1})
	AddEvent(w, evA{2})

	got, ok := GetEvents[evA](w)
	if !ok {
		t.Fatal("expected ok=true on non empty read")
	}
	if len(got) != 2 || got[0].v != 1 || got[1].v != 2 {
		t.Fatalf("unexpected first read: %#v", got)
	}

	got, ok = GetEvents[evA](w)
	if ok {
		t.Fatal("expected ok=false on empty read")
	}
	if len(got) != 0 {
		t.Fatalf("expected empty slice; got %#v", got)
	}
}

func TestDoubleBufferingAcrossUpdate(t *testing.T) {
	w := newTestWorld()
	setSystemInfo(w, 1, "sys")

	AddEvent(w, evA{10})
	update(w)

	AddEvent(w, evA{20})
	AddEvent(w, evA{30})

	got, ok := GetEvents[evA](w)
	if !ok {
		t.Fatal("expected ok=true")
	}
	if len(got) != 3 || got[0].v != 10 || got[1].v != 20 || got[2].v != 30 {
		t.Fatalf("unexpected read across update: %#v", got)
	}

	got, ok = GetEvents[evA](w)
	if ok {
		t.Fatal("expected ok=false after consuming")
	}
	if len(got) != 0 {
		t.Fatalf("expected empty after consuming; got %#v", got)
	}
}

func TestPerSystemIsolation(t *testing.T) {
	w := newTestWorld()

	setSystemInfo(w, 1, "sys1")
	AddEvent(w, evA{1})
	got1, ok := GetEvents[evA](w)
	if !ok || len(got1) != 1 || got1[0].v != 1 {
		t.Fatalf("sys1 wrong first read: %#v (ok=%v)", got1, ok)
	}

	setSystemInfo(w, 2, "sys2")
	got2, ok := GetEvents[evA](w)
	if !ok || len(got2) != 1 || got2[0].v != 1 {
		t.Fatalf("sys2 should independently see event: %#v (ok=%v)", got2, ok)
	}

	setSystemInfo(w, 1, "sys1")
	AddEvent(w, evA{2})
	got1, ok = GetEvents[evA](w)
	if !ok || len(got1) != 1 || got1[0].v != 2 {
		t.Fatalf("sys1 should see only new event: %#v (ok=%v)", got1, ok)
	}

	setSystemInfo(w, 2, "sys2")
	got2, ok = GetEvents[evA](w)
	if !ok || len(got2) != 1 || got2[0].v != 2 {
		t.Fatalf("sys2 should now see the second event: %#v (ok=%v)", got2, ok)
	}
}

func TestPerTypeIsolationSameSystem(t *testing.T) {
	w := newTestWorld()
	setSystemInfo(w, 1, "sys")

	AddEvent(w, evA{7})
	AddEvent(w, evB{"x"})

	as, ok := GetEvents[evA](w)
	if !ok || len(as) != 1 || as[0].v != 7 {
		t.Fatalf("unexpected A read: %#v (ok=%v)", as, ok)
	}

	bs, ok := GetEvents[evB](w)
	if !ok || len(bs) != 1 || bs[0].s != "x" {
		t.Fatalf("unexpected B read: %#v (ok=%v)", bs, ok)
	}
}

func TestUnknownTypeReturnsOkFalse(t *testing.T) {
	w := newTestWorld()
	setSystemInfo(w, 1, "sys")

	got, ok := GetEvents[evA](w)
	if ok {
		t.Fatal("expected ok=false when event type has never been created")
	}
	if len(got) != 0 {
		t.Fatalf("expected empty; got %#v", got)
	}
}

func TestSystemIdZeroIsUsable(t *testing.T) {
	w := newTestWorld()

	// Default systemInfo in a fresh world is id=0.
	AddEvent(w, evA{1})

	got, ok := GetEvents[evA](w)
	if !ok || len(got) != 1 || got[0].v != 1 {
		t.Fatalf("unexpected read with id=0: %#v (ok=%v)", got, ok)
	}
}

func TestConcurrentAddAndUpdateDoesNotPanic(t *testing.T) {
	w := newTestWorld()
	setSystemInfo(w, 1, "sys")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 10_000; i++ {
			AddEvent(w, evA{i})
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 2_000; i++ {
			update(w)
		}
	}()

	wg.Wait()

	// We do not assert exact counts here. This test is about safety (panic free).
}
