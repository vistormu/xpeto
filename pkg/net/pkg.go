package net

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/pkg/net/codec"
	"github.com/vistormu/xpeto/pkg/net/transport"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// resources
	ecs.AddResource(w, cache{
		packets:      make([]transport.Packet, 1024),
		packetCodecs: make([]codec.Codec, 1024),
	})
	ecs.AddResource(w, session{
		lookup:   make(map[string]ecs.Entity),
		channels: make([]Channel, 0),
	})

	// systems
	schedule.AddSystem(sch, schedule.PreUpdate, logTransportErrors,
		schedule.SystemOpt.Label("net.logTransportErrors"),
	)
	schedule.AddSystem(sch, schedule.PreUpdate, updateSession,
		schedule.SystemOpt.Label("net.updateSession"),
	)
	schedule.AddSystem(sch, schedule.PreUpdate, dispatch,
		schedule.SystemOpt.Label("net.dispatch"),
	)
	schedule.AddSystem(sch, schedule.PostUpdate, emitEvents,
		schedule.SystemOpt.Label("net.emitEvents"),
	)
}
