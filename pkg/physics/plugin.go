package physics

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/schedule"
)

func Plugin(ctx *core.Context, sch *schedule.Scheduler) {
	// resources
	core.AddResource(ctx, &PhysicsSettings{})

	// systems
	schedule.AddSystem(sch, schedule.FixedUpdate, applyGravity)
	schedule.AddSystem(sch, schedule.FixedUpdate, integrateVelocities)
	schedule.AddSystem(sch, schedule.FixedUpdate, buildBroadPhase)
	schedule.AddSystem(sch, schedule.FixedUpdate, narrowPhaseAABB)
	schedule.AddSystem(sch, schedule.FixedUpdate, resolveContactsAABB)
}
