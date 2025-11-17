package log

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// resources
	l := newLogger()
	l.sinks = append(l.sinks, &debugSink{})
	ecs.AddResource(w, l)

	// systems
	schedule.AddSystem(sch, schedule.Last, flush).Label("log.flush")
}
