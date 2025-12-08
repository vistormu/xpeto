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

// ===
// API
// ===
func SetTPS(w *ecs.World, tps int) {
	cs, _ := ecs.GetResource[ClockSettings](w)
	cs.FixedDelta = time.Second / time.Duration(tps)
}

func SetFixedDelta(w *ecs.World, s int) {
	cs, _ := ecs.GetResource[ClockSettings](w)
	cs.FixedDelta = time.Second * time.Duration(s)
}

func PauseClock(w *ecs.World, v bool) {
	cs, _ := ecs.GetResource[ClockSettings](w)
	cs.Paused = v
}
