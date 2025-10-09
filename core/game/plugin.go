package game

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

type Plugin func(*ecs.World, *schedule.Scheduler)
