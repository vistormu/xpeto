// app_test.go
package app

import (
	"errors"
	"testing"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

type testRes struct {
	V int
}

type testBackend struct {
	w   *ecs.World
	sch *schedule.Scheduler
	run func(*ecs.World, *schedule.Scheduler) error
}

func (b *testBackend) Run() error {
	if b.run == nil {
		return nil
	}
	return b.run(b.w, b.sch)
}

func TestRun_NilBackendFactory(t *testing.T) {
	a := NewApp(nil)
	err := a.Run()
	if err == nil {
		t.Fatalf("expected error for nil backend factory, got nil")
	}
}

func TestRun_StartupBeforeBackendRun_ExitAfterBackendRun(t *testing.T) {
	var gotW *ecs.World
	var gotSch *schedule.Scheduler

	backendFactory := func(w *ecs.World, sch *schedule.Scheduler) (Backend, error) {
		gotW = w
		gotSch = sch

		b := &testBackend{w: w, sch: sch}
		b.run = func(w *ecs.World, sch *schedule.Scheduler) error {
			// Startup must have already executed before backend.Run() is called.
			r, ok := ecs.GetResource[testRes](w)
			if !ok {
				return errors.New("missing testRes; pkg or startup system did not run")
			}
			if r.V != 1 {
				return errors.New("startup did not set testRes.V to 1 before backend.Run")
			}
			return nil
		}
		return b, nil
	}

	testPkg := func(w *ecs.World, sch *schedule.Scheduler) {
		// Make sure resource exists before systems execute.
		ecs.AddResource(w, testRes{V: 0})

		// Mark startup execution.
		schedule.AddSystem(sch, schedule.Startup, func(w *ecs.World) {
			r, _ := ecs.GetResource[testRes](w)
			r.V = 1
		}, schedule.SystemOpt.Label("test.startup"))

		// Mark exit execution.
		schedule.AddSystem(sch, schedule.Exit, func(w *ecs.World) {
			r, _ := ecs.GetResource[testRes](w)
			r.V = 2
		}, schedule.SystemOpt.Label("test.exit"))
	}

	a := NewApp(backendFactory, AppOpt.Pkgs(testPkg))

	err := a.Run()
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if gotW == nil || gotSch == nil {
		t.Fatalf("expected backendFactory to receive world and scheduler pointers")
	}

	// Exit runs via defer inside App.Run, so after Run returns the exit system must have executed.
	r, ok := ecs.GetResource[testRes](gotW)
	if !ok {
		t.Fatalf("expected testRes to exist after Run")
	}
	if r.V != 2 {
		t.Fatalf("expected testRes.V == 2 after exit, got %d", r.V)
	}
}

func TestRun_PkgNilIsIgnored(t *testing.T) {
	called := false

	backendFactory := func(w *ecs.World, sch *schedule.Scheduler) (Backend, error) {
		return &testBackend{
			w:   w,
			sch: sch,
			run: func(*ecs.World, *schedule.Scheduler) error { return nil },
		}, nil
	}

	p := func(*ecs.World, *schedule.Scheduler) { called = true }

	a := NewApp(backendFactory, AppOpt.Pkgs(nil, p, nil))

	err := a.Run()
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
	if !called {
		t.Fatalf("expected non-nil pkg to be called")
	}
}
