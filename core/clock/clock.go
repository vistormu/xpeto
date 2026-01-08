package clock

import (
	"math"
	"time"

	"github.com/vistormu/xpeto/core/ecs"
)

// =====
// clock
// =====
type NowFn func() time.Time

type RealClock struct {
	Start time.Time
	Now   time.Time
	Last  time.Time

	Delta   time.Duration
	Clamped time.Duration
	Elapsed time.Duration

	Frame uint64
}

func newRealClock(now time.Time) func() RealClock {
	return func() RealClock {
		return RealClock{
			Start: now,
			Now:   now,
			Last:  now,
		}
	}
}

type VirtualClock struct {
	Start   time.Time
	Delta   time.Duration
	Elapsed time.Duration

	Frame uint64
}

func newVirtualClock(now time.Time) func() VirtualClock {
	return func() VirtualClock {
		return VirtualClock{
			Start: now,
		}
	}
}

type FixedClock struct {
	Delta       time.Duration
	Accumulator time.Duration
	Steps       int

	Frame uint64
	Alpha float64
}

func newFixedClock() FixedClock {
	return FixedClock{
		Delta: time.Second / 60,
	}
}

// =======
// helpers
// =======
func scaledDelta(d time.Duration, scale float64, paused bool) time.Duration {
	if paused {
		return 0
	}
	if scale <= 0 {
		return 0
	}
	if math.IsNaN(scale) || math.IsInf(scale, 0) {
		scale = 1
	}

	x := float64(d) * scale
	if x <= 0 {
		return 0
	}

	if x > float64(math.MaxInt64) {
		return time.Duration(math.MaxInt64)
	}

	return time.Duration(x)
}

// =======
// systems
// =======
func tick(w *ecs.World) {
	s := ecs.EnsureResource(w, newClockSettings)

	if s.Now == nil {
		s.Now = time.Now
	}
	now := s.Now()

	real := ecs.EnsureResource(w, newRealClock(now))
	virtual := ecs.EnsureResource(w, newVirtualClock(now))
	fixed := ecs.EnsureResource(w, newFixedClock)

	// real clock
	realDelta := now.Sub(real.Last)
	real.Last = now
	real.Now = now

	if realDelta < 0 {
		realDelta = 0
	}
	real.Delta = realDelta
	real.Elapsed += realDelta
	real.Frame++

	clamped := realDelta
	if s.MaxDelta > 0 && clamped > s.MaxDelta {
		clamped = s.MaxDelta
	}
	real.Clamped = clamped

	// fixed delta intent
	if s.Mode == ModeFixed {
		fixed.Delta = max(0, s.FixedDelta)
	}

	// virtual clock
	virtualDelta := scaledDelta(clamped, s.Scale, s.Paused)
	if s.MaxVirtualDelta > 0 && virtualDelta > s.MaxVirtualDelta {
		virtualDelta = s.MaxVirtualDelta
	}

	virtual.Delta = virtualDelta
	virtual.Elapsed += virtualDelta
	virtual.Frame++

	// fixed clock
	fixed.Accumulator += virtualDelta

	maxSteps := s.MaxSteps
	if maxSteps <= 0 {
		maxSteps = 1
	}

	steps := 0
	if fixed.Delta > 0 {
		for fixed.Accumulator >= fixed.Delta && steps < maxSteps {
			fixed.Accumulator -= fixed.Delta
			steps++
			fixed.Frame++
		}
	}
	fixed.Steps = steps

	// alpha
	if fixed.Delta <= 0 {
		fixed.Alpha = 0
	} else {
		a := float64(fixed.Accumulator) / float64(fixed.Delta)
		if a < 0 {
			a = 0
		}
		if a >= 1 {
			a = math.Mod(a, 1)
		}
		fixed.Alpha = a
	}
}
