package render

import (
	"testing"

	"github.com/vistormu/xpeto/core"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

type testCanvas struct {
	calls []string
}

type opaqueItem struct {
	id  string
	key uint64
}

type transparentItem struct {
	id  string
	key uint64
}

func newTestWorld(t *testing.T) (*ecs.World, *schedule.Scheduler) {
	t.Helper()

	w := ecs.NewWorld()
	sch := schedule.NewScheduler()

	// Core must be installed so log/schedule resources exist.
	core.CorePkgs(w, sch)

	// Install render plugin under test.
	Pkg[testCanvas](w, sch)

	// Add canvas resource used by render().
	ok := ecs.AddResource(w, testCanvas{})
	if !ok {
		t.Fatalf("failed to add testCanvas resource (resources must be non-pointer types)")
	}

	return w, sch
}

func TestRender_DoesNotPanic_WhenResourcesMissing(t *testing.T) {
	t.Parallel()

	w := ecs.NewWorld()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("render panicked with missing resources: %v", r)
		}
	}()

	// Should log and return, not panic.
	render[testCanvas](w)
}

func TestPkg_RegistersRendererResource(t *testing.T) {
	t.Parallel()

	w, _ := newTestWorld(t)

	// This catches the common failure mode:
	// - Pkg adds a pointer resource (which AddResource rejects),
	// - renderer resource is never registered.
	_, ok := ecs.GetResource[renderer[testCanvas]](w)
	if !ok {
		t.Fatalf("expected renderer[testCanvas] resource to exist after Pkg; check that newRenderer returns a NON-pointer type")
	}
}

func TestAddRenderFn_DoesNotLeakBatch_WhenPrereqsMissing(t *testing.T) {
	t.Parallel()

	w, _ := newTestWorld(t)

	r, ok := ecs.GetResource[renderer[testCanvas]](w)
	if !ok || r == nil {
		t.Fatalf("expected renderer resource to exist")
	}

	before := len(r.batches)

	// Missing extraction + sort functions: registration must fail without changing state.
	AddRenderFn(w, Opaque, func(c *testCanvas, v opaqueItem) {
		c.calls = append(c.calls, v.id)
	})

	after := len(r.batches)
	if after != before {
		t.Fatalf("expected no new batch when prereqs are missing; before=%d after=%d", before, after)
	}
}

func TestRender_SortsByKey_AndRespectsStageOrder(t *testing.T) {
	t.Parallel()

	w, sch := newTestWorld(t)

	opaque := []opaqueItem{
		{id: "opaque-20", key: 20},
		{id: "opaque-10", key: 10},
	}
	trans := []transparentItem{
		{id: "trans-02", key: 2},
		{id: "trans-01", key: 1},
	}

	// Opaque pipeline
	AddExtractionFn[testCanvas](w, func(*ecs.World) []opaqueItem { return opaque })
	AddSortFn[testCanvas](w, func(v opaqueItem) uint64 { return v.key })
	AddRenderFn(w, Opaque, func(c *testCanvas, v opaqueItem) {
		c.calls = append(c.calls, v.id)
	})

	// Transparent pipeline
	AddExtractionFn[testCanvas](w, func(*ecs.World) []transparentItem { return trans })
	AddSortFn[testCanvas](w, func(v transparentItem) uint64 { return v.key })
	AddRenderFn(w, Transparent, func(c *testCanvas, v transparentItem) {
		c.calls = append(c.calls, v.id)
	})

	schedule.RunDraw(w, sch)

	canvas, _ := ecs.GetResource[testCanvas](w)
	got := canvas.calls

	// Expected default stage order for v0.1.0:
	// Opaque -> Transparent -> Ui -> PostFx
	want := []string{
		"opaque-10",
		"opaque-20",
		"trans-01",
		"trans-02",
	}

	if len(got) != len(want) {
		t.Fatalf("unexpected number of draw calls: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected draw order at %d: got=%v want=%v", i, got, want)
		}
	}
}

func TestRender_StableSort_WhenKeysEqual(t *testing.T) {
	t.Parallel()

	w, sch := newTestWorld(t)

	items := []opaqueItem{
		{id: "a", key: 1},
		{id: "b", key: 1},
		{id: "c", key: 1},
	}

	AddExtractionFn[testCanvas](w, func(*ecs.World) []opaqueItem { return items })
	AddSortFn[testCanvas](w, func(v opaqueItem) uint64 { return v.key })
	AddRenderFn(w, Opaque, func(c *testCanvas, v opaqueItem) {
		c.calls = append(c.calls, v.id)
	})

	schedule.RunDraw(w, sch)

	canvas, _ := ecs.GetResource[testCanvas](w)
	got := canvas.calls
	want := []string{"a", "b", "c"}

	if len(got) != len(want) {
		t.Fatalf("unexpected number of draw calls: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("stable order violated at %d: got=%v want=%v", i, got, want)
		}
	}
}
