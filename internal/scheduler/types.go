package scheduler

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"
)

type Schedule struct {
	Name      string
	Stage     core.Stage
	System    ecs.System
	Before    []string
	After     []string
	Condition func(*core.Context) bool
}
