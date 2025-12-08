package time

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	xptime "github.com/vistormu/xpeto/core/time"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	s, _ := ecs.GetResource[xptime.ClockSettings](w)
	ecs.AddResource(w, lastClockSettings{*s})

	schedule.AddSystem(sch, schedule.PreUpdate, watch)
}
