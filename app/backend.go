package app

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

type BackendFactory = func(*ecs.World, *schedule.Scheduler) (Backend, error)

type Backend interface {
	Run() error
}
