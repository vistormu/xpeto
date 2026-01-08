package clock

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/core/clock"
	"github.com/vistormu/xpeto/core/ecs"
)

type lastClockSettings struct {
	Mode       clock.ClockMode
	FixedDelta time.Duration
}

func newLastClockSettings(s *clock.ClockSettings) func() lastClockSettings {
	return func() lastClockSettings {
		return lastClockSettings{
			Mode:       s.Mode,
			FixedDelta: s.FixedDelta,
		}
	}
}

func watch(w *ecs.World) {
	s, ok := ecs.GetResource[clock.ClockSettings](w)
	if !ok {
		return
	}
	ls := ecs.EnsureResource(w, newLastClockSettings(s))

	// mode
	if s.Mode != ls.Mode {
		switch s.Mode {
		case clock.ModeFixed:
			ebiten.SetTPS(s.Tps)
		case clock.ModeSyncWithFPS:
			ebiten.SetTPS(ebiten.SyncWithFPS)
		default:
			ebiten.SetTPS(ebiten.DefaultTPS)
		}

		ls.Mode = s.Mode
	}

	// tps and fixed delta
	if (s.Mode == clock.ModeFixed) && (s.FixedDelta != ls.FixedDelta) {
		ebiten.SetTPS(s.Tps)
		ls.FixedDelta = s.FixedDelta
	}
}
