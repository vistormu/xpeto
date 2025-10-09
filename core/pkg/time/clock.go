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
}

// =====
// clock
// =====
type Clock struct {
	Start     time.Time
	LastFrame time.Time
	Elapsed   time.Duration
	Delta     time.Duration
	Frame     uint64
}

type fixedClock struct {
	maxFixedSteps int
	accumulator   time.Duration
	fixedSteps    int
}

// =======
// systems
// =======
func tickClockSystem(w *ecs.World) {
	clk, ok := ecs.GetResource[Clock](w)
	if !ok || clk == nil {
		return
	}
	fc, ok := ecs.GetResource[fixedClock](w)
	if !ok || fc == nil {
		return
	}

	// Defaults (used if no ClockSettings present)
	settings := ClockSettings{
		FixedDelta:  stdtime.Second / 60, // 60 Hz fixed
		Scale:       1.0,
		Paused:      false,
		SyncWithFps: false,
	}
	if cs, ok := ecs.GetResource[ClockSettings](w); ok && cs != nil {
		settings = *cs
	}

	// Hand off sync to Ebitengine (idempotent; cheap if unchanged).
	ebiten.SetVsyncEnabled(settings.SyncWithFps)

	// First-tick init
	if clk.Start.IsZero() {
		now := stdtime.Now()
		clk.Start = now
		clk.LastFrame = now
	}

	// Real delta from wall clock
	now := stdtime.Now()
	realDelta := now.Sub(clk.LastFrame)
	clk.LastFrame = now

	// Variable (virtual) delta = real * scale (unless paused)
	var vdt stdtime.Duration
	if settings.Paused || settings.Scale == 0 {
		vdt = 0
	} else {
		vdt = stdtime.Duration(float64(realDelta) * settings.Scale)
	}

	// Update frame clock
	clk.Delta = vdt
	clk.Elapsed += vdt
	clk.Frame++

	// --- Fixed-step accumulation ---
	if fc.maxFixedSteps <= 0 {
		fc.maxFixedSteps = 8 // avoid spiral-of-death
	}

	fd := settings.FixedDelta
	fc.accumulator += vdt

	steps := 0
	for fd > 0 && fc.accumulator >= fd && steps < fc.maxFixedSteps {
		fc.accumulator -= fd
		steps++
	}
	fc.fixedSteps = steps
}

// Helpers your scheduler/systems can call:

// FixedStepsThisFrame returns how many times to run the fixed stages this frame.
func FixedStepsThisFrame(w *ecs.World) int {
	if fc, ok := ecs.GetResource[fixedClock](w); ok && fc != nil {
		return fc.fixedSteps
	}
	return 0
}

// FixedDelta returns the constant dt each fixed tick should use.
func FixedDelta(w *ecs.World) stdtime.Duration {
	if cs, ok := ecs.GetResource[ClockSettings](w); ok && cs != nil && cs.FixedDelta > 0 {
		return cs.FixedDelta
	}
	return stdtime.Second / 60
}
