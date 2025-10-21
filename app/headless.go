//go:build headless

package app

import (
	"time"

	"github.com/vistormu/go-dsa/system"
	"github.com/vistormu/xpeto/core/ecs"
	xptime "github.com/vistormu/xpeto/core/pkg/time"
)

var toRunner = map[Runner]runner{
	Ebiten:   &headlessRunner{},
	Headless: &headlessRunner{},
}

type headlessRunner struct {
	app *App
}

func (r *headlessRunner) run(a *App) error {
	a.build()
	r.app = a
	r.update()
	return nil
}

func (r *headlessRunner) update() {
	stopper := system.NewKbIntListener()
	defer stopper.Stop()

	cs, _ := ecs.GetResource[xptime.ClockSettings](r.app.world)

	timer := time.NewTimer(cs.FixedDelta)
	defer timer.Stop()

	for {
		select {
		case <-stopper.Listen():
			return

		case <-timer.C:
			r.app.scheduler.RunUpdate(r.app.world)

			latest, _ := ecs.GetResource[xptime.ClockSettings](r.app.world)

			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}

			timer.Reset(latest.FixedDelta)
		}
	}
}
