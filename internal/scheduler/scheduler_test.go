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

	stages := StartupStages()
	stages = append(stages, UpdateStages()...)

	scheduler := NewScheduler(stages).
		WithSchedule(&Schedule{Name: "preStartup", Stage: PreStartup, System: preStartupSystem}).
		WithSchedule(&Schedule{Name: "startup", Stage: Startup, System: startupSystem}).
		WithSchedule(&Schedule{Name: "postStartup", Stage: PostStartup, System: postStartupSystem}).
		WithSchedule(&Schedule{Name: "first", Stage: First, System: firstSystem}).
		WithSchedule(&Schedule{Name: "preUpdate", Stage: PreUpdate, System: preUpdateSystem}).
		WithSchedule(&Schedule{Name: "fixedFirst", Stage: FixedFirst, System: fixedFirstSystem}).
		WithSchedule(&Schedule{Name: "fixedPreUpdate", Stage: FixedPreUpdate, System: fixedPreUpdateSystem}).
		WithSchedule(&Schedule{Name: "fixedUpdate", Stage: FixedUpdate, System: fixedUpdateSystem}).
		WithSchedule(&Schedule{Name: "fixedPostUpdate", Stage: FixedPostUpdate, System: fixedPostUpdateSystem}).
		WithSchedule(&Schedule{Name: "fixedLast", Stage: FixedLast, System: fixedLastSystem}).
		WithSchedule(&Schedule{Name: "update", Stage: Update, System: updateSystem}).
		WithSchedule(&Schedule{Name: "postUpdate", Stage: PostUpdate, System: postUpdateSystem}).
		WithSchedule(&Schedule{Name: "last", Stage: Last, System: lastSystem})

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
