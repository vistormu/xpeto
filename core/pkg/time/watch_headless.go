//go:build headless

package time

import (
	"time"

	"github.com/vistormu/xpeto/core/ecs"
)

func applyChanges(w *ecs.World) {
	fixed, _ := ecs.GetResource[FixedClock](w)
	s, _ := ecs.GetResource[ClockSettings](w)
	ls, _ := ecs.GetResource[lastClockSettings](w)

	if s.FixedDelta != ls.FixedDelta {
		if s.FixedDelta <= 0 {
			fixed.Delta = time.Second / 60
		} else {
			hz := 1.0 / s.FixedDelta.Seconds()
			if hz < 1 {
				fixed.Delta = time.Second
			} else if hz > 10_000 {
				fixed.Delta = time.Second / 10_000
			} else {
				fixed.Delta = s.FixedDelta
			}
		}
		ls.FixedDelta = s.FixedDelta
	}

	ls.Scale = s.Scale
	ls.Paused = s.Paused
	ls.MaxDelta = s.MaxDelta
}
