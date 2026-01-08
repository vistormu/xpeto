package schedule

import (
	"strings"
	"testing"

	"github.com/vistormu/xpeto/core/ecs"
)

func TestSchedule_OrderWithLabels(t *testing.T) {
	w := ecs.NewWorld()
	sch := NewScheduler()

	order := make([]string, 0)
	push := func(s string) ecs.System {
		return func(*ecs.World) { order = append(order, s) }
	}

	AddSystem(sch, Update, push("A"),
		SystemOpt.Label("A"),
	)
	AddSystem(sch, Update, push("B"),
		SystemOpt.Label("B"),
		SystemOpt.RunBefore("C"),
	)
	AddSystem(sch, Update, push("C"),
		SystemOpt.Label("C"),
		SystemOpt.RunAfter("A"),
	)

	RunStartup(w, sch)
	RunUpdate(w, sch)

	got := strings.Join(order, "")
	if got != "ABC" {
		t.Fatalf("expected order ABC, got %q", got)
	}
}

// func TestSchedule_UnknownLabelDiagnostic(t *testing.T) {
// 	w := ecs.NewWorld()
// 	sch := NewScheduler()

// 	var diags []Diagnostic
// 	SetDiagnosticSink(sch, func(_ *ecs.World, d Diagnostic) { diags = append(diags, d) })

// 	AddSystem(sch, Update, func(*ecs.World) {},
// 		SystemOpt.Label("A"),
// 	)
// 	AddSystem(sch, Update, func(*ecs.World) {},
// 		SystemOpt.RunAfter("missing"),
// 	)

// 	RunStartup(w, sch)
// 	RunUpdate(w, sch)

// 	found := false
// 	for _, d := range diags {
// 		if strings.Contains(d.Message, "unknown dependency label") {
// 			found = true
// 			break
// 		}
// 	}
// 	if !found {
// 		t.Fatalf("expected an unknown-label diagnostic, got %#v", diags)
// 	}
// }

// func TestSchedule_CycleFallsBackToInsertionOrder(t *testing.T) {
// 	w := ecs.NewWorld()
// 	sch := NewScheduler()

// 	var diags []Diagnostic
// 	SetDiagnosticSink(sch, func(_ *ecs.World, d Diagnostic) { diags = append(diags, d) })

// 	order := make([]string, 0)
// 	AddSystem(sch, Update, func(*ecs.World) { order = append(order, "A") },
// 		SystemOpt.Label("A"),
// 		SystemOpt.RunAfter("B"),
// 	)
// 	AddSystem(sch, Update, func(*ecs.World) { order = append(order, "B") },
// 		SystemOpt.Label("B"),
// 		SystemOpt.RunAfter("A"),
// 	)

// 	RunStartup(w, sch)
// 	RunUpdate(w, sch)

// 	got := strings.Join(order, "")
// 	if got != "AB" {
// 		t.Fatalf("expected insertion order AB on cycle, got %q", got)
// 	}

// 	found := false
// 	for _, d := range diags {
// 		if strings.Contains(d.Message, "cycle detected") {
// 			found = true
// 			break
// 		}
// 	}
// 	if !found {
// 		t.Fatalf("expected a cycle diagnostic, got %#v", diags)
// 	}
// }

func TestSchedule_ConditionOnce(t *testing.T) {
	w := ecs.NewWorld()
	sch := NewScheduler()

	count := 0
	AddSystem(sch, Update, func(*ecs.World) { count++ },
		SystemOpt.RunIf(Once()),
	)

	RunStartup(w, sch)
	RunUpdate(w, sch)
	RunUpdate(w, sch)

	if count != 1 {
		t.Fatalf("expected system to run once, got count=%d", count)
	}
}

func TestSchedule_StateCallbacksExecute(t *testing.T) {
	type S int
	const (
		Idle S = iota
		Running
	)

	w := ecs.NewWorld()
	sch := NewScheduler()
	AddStateMachine(sch, Idle)

	log := make([]string, 0)
	AddSystem(sch, OnEnter(Running), func(*ecs.World) { log = append(log, "enter") })
	AddSystem(sch, OnExit(Idle), func(*ecs.World) { log = append(log, "exit") })
	AddSystem(sch, OnTransition(Idle, Running), func(*ecs.World) { log = append(log, "transition") })

	RunStartup(w, sch)

	if !SetNextState(w, Running) {
		t.Fatalf("expected SetNextState to succeed after startup")
	}
	RunUpdate(w, sch)

	got := strings.Join(log, ",")
	if got != "exit,transition,enter" {
		t.Fatalf("expected callbacks exit,transition,enter; got %q", got)
	}
}
