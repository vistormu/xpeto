package scheduler

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"
)

type Stage uint32

const (
	PreStartup Stage = iota
	Startup
	PostStartup

	First
	PreUpdate

	FixedFirst
	FixedPreUpdate
	FixedUpdate
	FixedPostUpdate
	FixedLast

	Update
	PostUpdate
	Last
)

func StartupStages() []Stage {
	return []Stage{
		PreStartup,
		Startup,
		PostStartup,
	}
}

func UpdateStages() []Stage {
	return []Stage{
		First,
		PreUpdate,
		FixedFirst,
		FixedPreUpdate,
		FixedUpdate,
		FixedPostUpdate,
		FixedLast,
		Update,
		PostUpdate,
		Last,
	}
}

type Schedule struct {
	Name      string
	Stage     Stage
	System    ecs.System
	Before    []string
	After     []string
	Condition func(*core.Context) bool
}
