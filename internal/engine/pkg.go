package engine

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/scheduler"
)

type Pkg interface {
	Resources() []any
	Schedules() []*scheduler.Schedule
	Build(*core.Context)
}
