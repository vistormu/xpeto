package physics

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// resources
	ecs.AddResource(w, Space{
		Cells:      make([]*Cell, 0),
		candidates: make([]pair, 0),
		Contacts:   make([]ContactPair, 0),
	})
	ecs.AddResource(w, lastSpaceSize{})
	ecs.AddResource(w, Gravity{})

	// systems
	schedule.AddSystem(sch, schedule.FixedPreUpdate, resizeSpace)

	schedule.AddSystem(sch, schedule.FixedUpdate, applyGravity)
	schedule.AddSystem(sch, schedule.FixedUpdate, integrateVelocities)

	schedule.AddSystem(sch, schedule.FixedUpdate, fillGrid)
	schedule.AddSystem(sch, schedule.FixedUpdate, getCandidates)

	schedule.AddSystem(sch, schedule.FixedUpdate, narrowPhaseRectRect)
	schedule.AddSystem(sch, schedule.FixedUpdate, resolveContacts)
	schedule.AddSystem(sch, schedule.FixedUpdate, resolveContactsImpulses)
}
