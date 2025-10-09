package time

import (
	"time"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func Pkg(w *ecs.World, sch *schedule.Scheduler) {
	// resources
	ecs.AddResource(w, clockSettings{
		fixedDelta:  time.Second / 60,
		scale:       1.0,
		paused:      false,
		syncWithFps: false,
	})
	ecs.AddResource(w, clock{
		start:     time.Now(),
		lastFrame: time.Now(),
		elapsed:   0,
		delta:     0,
		frame:     0,
	})
	ecs.AddResource(w, fixedClock{
		maxFixedSteps: 8,
		accumulator:   0,
		fixedSteps:    0,
	})
}
