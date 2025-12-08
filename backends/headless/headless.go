package headless

import (
	"time"

	"github.com/vistormu/go-dsa/system"
	"github.com/vistormu/xpeto/app"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/event"
	"github.com/vistormu/xpeto/core/schedule"
	xptime "github.com/vistormu/xpeto/core/time"
)

// ===============
// headless runner
// ===============
type runner struct {
	w   *ecs.World
	sch *schedule.Scheduler
}

func (r *runner) update() {
	stopper := system.NewKbIntListener()
	defer stopper.Stop()

	cs, _ := ecs.GetResource[xptime.ClockSettings](r.w)

	timer := time.NewTimer(cs.FixedDelta)
	defer timer.Stop()

	for {
		select {
		case <-stopper.Listen():
			return

		case <-timer.C:
			r.sch.RunUpdate(r.w)

			// dynamic tick
			latest, _ := ecs.GetResource[xptime.ClockSettings](r.w)

			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}

			timer.Reset(latest.FixedDelta)

			// exit event
			_, ok := event.GetEvents[app.ExitApp](r.w)
			if ok {
				return
			}
		}
	}
}

// =======
// backend
// =======
func Backend() app.Backend {
	return &backend{}
}

type backend struct {
	w   *ecs.World
	sch *schedule.Scheduler
}

func (b *backend) Init(w *ecs.World, sch *schedule.Scheduler) {
	b.w = w
	b.sch = sch
	corePkgs(w, sch)
}

func (b *backend) Run() error {
	r := new(runner)
	r.w = b.w
	r.sch = b.sch

	r.update()

	return nil
}
