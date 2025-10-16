//go:build headless
// +build headless

package app

import (
	"time"

	"github.com/vistormu/go-dsa/system"

	"github.com/vistormu/xpeto/core/ecs"
	xptime "github.com/vistormu/xpeto/core/pkg/time"
)

type headlessRunner struct {
	app *App
}

func (r *headlessRunner) update() {
	stopper := system.NewKbIntListener()
	defer stopper.Stop()

	cs, _ := ecs.GetResource[xptime.ClockSettings](r.app.world)

	ticker := time.NewTicker(cs.FixedDelta)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-stopper.Listen():
			break loop

		case <-ticker.C:
			r.app.scheduler.RunUpdate(r.app.world)
		}
	}
}

func (r *headlessRunner) run(a *App) error {
	a.build()

	r.app = a

	r.update()

	return nil
}
