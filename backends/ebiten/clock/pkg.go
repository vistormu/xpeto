package clock

import (
	"github.com/vistormu/xpeto/core/clock"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/log"
	"github.com/vistormu/xpeto/core/schedule"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	s, ok := ecs.GetResource[clock.ClockSettings](w)
	if !ok {
		log.LogError(w, "missing clock.ClockSettings for ebiten.clock.Pkg")
		return
	}
	ecs.AddResource(w, newLastClockSettings(s)())

	schedule.AddSystem(sch, schedule.First, watch,
		schedule.SystemOpt.Label("ebiten.clock.watch"),
		schedule.SystemOpt.RunAfter("clock.sanitizeClockSettings"),
		schedule.SystemOpt.RunBefore("clock.tick"),
	)
}
