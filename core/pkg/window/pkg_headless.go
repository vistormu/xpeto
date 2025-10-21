//go:build headless

package window

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

type Window struct{}

func Pkg(w *ecs.World, sch *schedule.Scheduler) {}
