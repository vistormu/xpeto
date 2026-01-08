package core

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/core/window"

	"github.com/vistormu/xpeto/core/clock"
	"github.com/vistormu/xpeto/core/event"
	"github.com/vistormu/xpeto/core/log"
	"github.com/vistormu/xpeto/core/log/sink"
)

type Pkg = func(*ecs.World, *schedule.Scheduler)

func CorePkgs(w *ecs.World, sch *schedule.Scheduler) {
	event.Pkg(w, sch)
	clock.Pkg(w, sch)
	window.Pkg(w, sch)
	log.Pkg(w, sch)
	log.AddSink(w, &sink.TerminalSink{})
}
