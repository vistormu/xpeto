package log

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// resources
	s := newLoggerSettings()
	l := newLogger(s)()
	ecs.AddResource(w, l)
	ecs.AddResource(w, s)

	// systems
	schedule.AddSystem(sch, schedule.Last, flush,
		schedule.SystemOpt.Label("log.flush"),
	)
	schedule.AddSystem(sch, schedule.Exit, flush,
		schedule.SystemOpt.Label("log.flush_exit"),
	)
}
