package time

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/core/ecs"
	xptime "github.com/vistormu/xpeto/core/time"
)

type lastClockSettings struct {
	xptime.ClockSettings
}

func watch(w *ecs.World) {
	fixed, _ := ecs.GetResource[xptime.FixedClock](w)
	s, _ := ecs.GetResource[xptime.ClockSettings](w)
	ls, _ := ecs.GetResource[lastClockSettings](w)

	// sync with fps
	if s.SyncWithFps != ls.SyncWithFps {
		if s.SyncWithFps {
			ebiten.SetTPS(ebiten.SyncWithFPS)
		}
		ls.SyncWithFps = s.SyncWithFps
	}

	// fixed delta
	if !s.SyncWithFps && s.FixedDelta != ls.FixedDelta {
		if s.FixedDelta <= 0 {
			ebiten.SetTPS(ebiten.DefaultTPS)
			fixed.Delta = time.Second / time.Duration(ebiten.DefaultTPS)
		} else {
			hz := 1.0 / s.FixedDelta.Seconds()
			tps := int(math.Round(hz))
			if tps < 1 {
				tps = 1
			} else if tps > 10_000 {
				tps = 10_000
			}
			ebiten.SetTPS(tps)
			fixed.Delta = s.FixedDelta
		}

		ls.FixedDelta = s.FixedDelta
	}

	ls.Scale = s.Scale
	ls.Paused = s.Paused
	ls.MaxDelta = s.MaxDelta
}
