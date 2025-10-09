package pkg

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

type Pkg = func(*ecs.World, *schedule.Scheduler)
