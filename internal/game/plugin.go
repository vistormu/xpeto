package game

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/schedule"
)

type Plugin func(*core.Context, *schedule.Scheduler)
