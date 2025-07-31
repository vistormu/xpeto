package scheduler

import (
	"testing"

	"github.com/vistormu/xpeto/internal/core"
)

// =====
// tests
// =====
func TestStages(t *testing.T) {
	ctx := core.NewContext()
	core.AddResource(ctx, &[]string{})

	stages := core.StartupStages()
	stages = append(stages, core.UpdateStages()...)

	scheduler := NewScheduler(stages).
		WithSchedule(&Schedule{Name: "preStartup", Stage: core.PreStartup, System: preStartupSystem}).
		WithSchedule(&Schedule{Name: "startup", Stage: core.Startup, System: startupSystem}).
		WithSchedule(&Schedule{Name: "postStartup", Stage: core.PostStartup, System: postStartupSystem}).
		WithSchedule(&Schedule{Name: "first", Stage: core.First, System: firstSystem}).
		WithSchedule(&Schedule{Name: "preUpdate", Stage: core.PreUpdate, System: preUpdateSystem}).
		WithSchedule(&Schedule{Name: "fixedFirst", Stage: core.FixedFirst, System: fixedFirstSystem}).
		WithSchedule(&Schedule{Name: "fixedPreUpdate", Stage: core.FixedPreUpdate, System: fixedPreUpdateSystem}).
		WithSchedule(&Schedule{Name: "fixedUpdate", Stage: core.FixedUpdate, System: fixedUpdateSystem}).
		WithSchedule(&Schedule{Name: "fixedPostUpdate", Stage: core.FixedPostUpdate, System: fixedPostUpdateSystem}).
		WithSchedule(&Schedule{Name: "fixedLast", Stage: core.FixedLast, System: fixedLastSystem}).
		WithSchedule(&Schedule{Name: "update", Stage: core.Update, System: updateSystem}).
		WithSchedule(&Schedule{Name: "postUpdate", Stage: core.PostUpdate, System: postUpdateSystem}).
		WithSchedule(&Schedule{Name: "last", Stage: core.Last, System: lastSystem})

	if len(scheduler.schedulesByStage) != 13 {
		t.Fatalf("expected 13 stages, got %d", len(scheduler.schedulesByStage))
	}

	scheduler.Run(ctx)

	order := core.MustResource[*[]string](ctx)
	expectedOrder := []string{
		"preStartup",
		"startup",
		"postStartup",
		"first",
		"preUpdate",
		"fixedFirst",
		"fixedPreUpdate",
		"fixedUpdate",
		"fixedPostUpdate",
		"fixedLast",
		"update",
		"postUpdate",
		"last",
	}

	if len(*order) != len(expectedOrder) {
		t.Fatalf("expected %d systems to run, got %d", len(expectedOrder), len(*order))
	}

	for i, sys := range *order {
		if sys != expectedOrder[i] {
			t.Errorf("expected system %d to be %s, got %s", i, expectedOrder[i], sys)
		}
	}
}

func TestEmptyScheduler(t *testing.T) {
	ctx := core.NewContext()
	scheduler := NewScheduler(core.StartupStages())

	scheduler.Run(ctx)
}

// =======
// helpers
// =======
func mutateOrder(ctx *core.Context, message string) {
	order := core.MustResource[*[]string](ctx)
	*order = append(*order, message)
}

func preStartupSystem(ctx *core.Context)      { mutateOrder(ctx, "preStartup") }
func startupSystem(ctx *core.Context)         { mutateOrder(ctx, "startup") }
func postStartupSystem(ctx *core.Context)     { mutateOrder(ctx, "postStartup") }
func firstSystem(ctx *core.Context)           { mutateOrder(ctx, "first") }
func preUpdateSystem(ctx *core.Context)       { mutateOrder(ctx, "preUpdate") }
func fixedFirstSystem(ctx *core.Context)      { mutateOrder(ctx, "fixedFirst") }
func fixedPreUpdateSystem(ctx *core.Context)  { mutateOrder(ctx, "fixedPreUpdate") }
func fixedUpdateSystem(ctx *core.Context)     { mutateOrder(ctx, "fixedUpdate") }
func fixedPostUpdateSystem(ctx *core.Context) { mutateOrder(ctx, "fixedPostUpdate") }
func fixedLastSystem(ctx *core.Context)       { mutateOrder(ctx, "fixedLast") }
func updateSystem(ctx *core.Context)          { mutateOrder(ctx, "update") }
func postUpdateSystem(ctx *core.Context)      { mutateOrder(ctx, "postUpdate") }
func lastSystem(ctx *core.Context)            { mutateOrder(ctx, "last") }
