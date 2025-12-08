package app

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

type Backend interface {
	Init(w *ecs.World, sch *schedule.Scheduler)
	Run() error
}
