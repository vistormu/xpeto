package time

import (
	"time"
)

// Time represents the time state of the game engine.
// It tracks the current frame, elapsed time, and allows for time scaling.
// - Delta: the time since the last frame
// - ScaledDelta: Delta adjusted by Scale factor unless paused
// - Elapsed: total virtual time since the game started
// - RealElapsed: total wall-clock time since the game started
// - LastReal: the last real timestamp for delta calculations
// - Frame: the total number of frames processed
// - Paused: whether the game is currently paused
// - Scale: a factor to scale time (1 = real-time, 0.5 = half-speed, etc.)
// - maxDelta: an optional maximum delta to clamp spikes in frame time
type Time struct {
	Delta       time.Duration
	ScaledDelta time.Duration
	Elapsed     time.Duration
	RealElapsed time.Duration
	LastReal    time.Time
	Frame       uint64
	Paused      bool
	Scale       float64
	maxDelta    time.Duration
}
