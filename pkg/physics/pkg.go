package physics

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// resources
	ecs.AddResource(w, &Settings{})

	// systems
	schedule.AddSystem(sch, schedule.FixedUpdate, applyGravity)
	schedule.AddSystem(sch, schedule.FixedUpdate, integrateVelocities)

	cs := NewCollisionSolver()
	schedule.AddSystem(sch, schedule.FixedUpdate, cs.buildBroadPhase)
	schedule.AddSystem(sch, schedule.FixedUpdate, cs.narrowPhaseAABB)
	schedule.AddSystem(sch, schedule.FixedUpdate, cs.resolveContactsAABB)
}
