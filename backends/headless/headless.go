package headless

import (
	"fmt"
	"runtime"
	"time"

	"github.com/vistormu/go-dsa/system"
	"github.com/vistormu/xpeto/app"
	"github.com/vistormu/xpeto/core/clock"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/event"
	"github.com/vistormu/xpeto/core/schedule"
)

// =======
// backend
// =======
func Backend(w *ecs.World, sch *schedule.Scheduler) (app.Backend, error) {
	runner := &runner{
		w:   w,
		sch: sch,
	}

	return runner, nil
}

// =======
// helpers
// =======
func resetTimer(t *time.Timer, d time.Duration) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
	t.Reset(d)
}

// ===============
// headless runner
// ===============
type runner struct {
	w   *ecs.World
	sch *schedule.Scheduler
}

func (r *runner) Run() error {
	stopper := system.NewSignalListener()
	defer stopper.Stop()

	stop := stopper.Listen()

	timer := time.NewTimer(0)
	if !timer.Stop() {
		<-timer.C
	}
	defer timer.Stop()

	for {
		cs, ok := ecs.GetResource[clock.ClockSettings](r.w)
		if !ok {
			return fmt.Errorf("missing clocl.ClockSettings resource. cannot compute loop")
		}

		delta := cs.FixedDelta
		if delta <= 0 {
			delta = time.Second / 60
		}

		switch cs.Mode {
		case clock.ModeFixed:
			fixed := cs.FixedDelta
			if fixed <= 0 {
				return fmt.Errorf("clock: FixedDelta must be > 0 in ModeFixed")
			}

			timer := time.NewTimer(0)
			defer timer.Stop()

			nextTick := time.Now()

			for {
				select {
				case <-stop:
					return nil
				default:
				}

				now := time.Now()
				wait := time.Until(nextTick)

				if wait > 0 {
					resetTimer(timer, wait)

					select {
					case <-stop:
						return nil
					case <-timer.C:
					}
				}

				schedule.RunUpdate(r.w, r.sch)

				if _, ok := event.GetEvents[app.ExitAppEvent](r.w); ok {
					return nil
				}

				nextTick = nextTick.Add(fixed)

				now = time.Now()
				if lag := now.Sub(nextTick); lag > fixed {
					missed := lag / fixed
					nextTick = nextTick.Add(missed * fixed)
				}
			}

		case clock.ModeSyncWithFPS:
			select {
			case <-stop:
				return nil
			default:
			}

			schedule.RunUpdate(r.w, r.sch)

			if _, ok := event.GetEvents[app.ExitAppEvent](r.w); ok {
				return nil
			}

			runtime.Gosched()

		default:
			cs.Mode = clock.ModeSyncWithFPS
		}
	}
}
