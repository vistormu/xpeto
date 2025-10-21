package time

import (
	"time"

	"github.com/vistormu/xpeto/core/ecs"
)

// ========
// settings
// ========
type ClockSettings struct {
	FixedDelta  time.Duration
	Scale       float64
	Paused      bool
	SyncWithFps bool
	MaxDelta    time.Duration
}

type lastClockSettings struct {
	ClockSettings
}

// =====
// clock
// =====
type RealClock struct {
	Start   time.Time
	Last    time.Time
	Delta   time.Duration
	Clamped time.Duration
	Elapsed time.Duration
}

type VirtualClock struct {
	Start   time.Time
	Delta   time.Duration
	Elapsed time.Duration
	Frame   uint64
}

type FixedClock struct {
	Delta       time.Duration
	Accumulator time.Duration
	Steps       int
	MaxSteps    int
}

// =======
// systems
// =======
func tick(w *ecs.World) {
	real, _ := ecs.GetResource[RealClock](w)
	virtual, _ := ecs.GetResource[VirtualClock](w)
	fixed, _ := ecs.GetResource[FixedClock](w)
	s, _ := ecs.GetResource[ClockSettings](w)

	now := time.Now()

	// real clock
	realDelta := now.Sub(real.Last)
	real.Last = now

	if realDelta < 0 {
		realDelta = 0
	}
	real.Delta = realDelta

	real.Elapsed += realDelta

	clamped := realDelta
	if s.MaxDelta > 0 && clamped > s.MaxDelta {
		clamped = s.MaxDelta
	}
	real.Clamped = clamped

	// virtual clock
	var virtualDelta time.Duration
	if s.Paused || s.Scale <= 0 {
		virtualDelta = 0
	} else {
		scaled := float64(clamped) * s.Scale
		virtualDelta = time.Duration(scaled)
	}

	virtual.Delta = virtualDelta
	virtual.Elapsed += virtualDelta

	virtual.Frame++

	// fixed clock
	fixed.Accumulator += virtualDelta

	steps := 0
	for fixed.Delta > 0 && fixed.Accumulator >= fixed.Delta && steps < fixed.MaxSteps {
		fixed.Accumulator -= fixed.Delta
		steps++
	}

	fixed.Steps = steps
}

func applyInitialSettings(w *ecs.World) {
	real, _ := ecs.GetResource[RealClock](w)
	virtual, _ := ecs.GetResource[VirtualClock](w)
	s, _ := ecs.GetResource[ClockSettings](w)

	// set initial time
	now := time.Now()
	if real.Start.IsZero() {
		real.Start = now
		real.Last = now
	}
	if virtual.Start.IsZero() {
		virtual.Start = now
	}

	// add internal resouce to track changes
	ecs.AddResource(w, lastClockSettings{*s})
}

// ===
// API
// ===
