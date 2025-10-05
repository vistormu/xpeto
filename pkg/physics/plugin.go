package physics

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/schedule"
)

func Plugin(ctx *core.Context, sch *schedule.Scheduler) {
	// resources
	core.AddResource(ctx, &Settings{})

	// systems
	schedule.AddSystem(sch, schedule.FixedUpdate, applyGravity)
	schedule.AddSystem(sch, schedule.FixedUpdate, integrateVelocities)

	cs := NewCollisionSolver()
	schedule.AddSystem(sch, schedule.FixedUpdate, cs.buildBroadPhase)
	schedule.AddSystem(sch, schedule.FixedUpdate, cs.narrowPhaseAABB)
	schedule.AddSystem(sch, schedule.FixedUpdate, cs.resolveContactsAABB)
}
