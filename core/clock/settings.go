package clock

import (
	"math"
	"time"

	"github.com/vistormu/xpeto/core/ecs"
)

type ClockMode uint8

const (
	ModeFixed ClockMode = iota
	ModeSyncWithFPS
)

const (
	DefaultFixedDelta      = time.Second / 60
	DefaultTps             = 60
	DefaultMaxDelta        = 100 * time.Millisecond
	DefaultMaxVirtualDelta = 0
	DefaultScale           = 1.0
	DefaultPaused          = false
	DefaultMaxSteps        = 8
	MinTps                 = 0
	MaxTps                 = 10_000
)

// ========
// settings
// ========
type ClockSettings struct {
	Mode ClockMode
	Now  NowFn

	FixedDelta time.Duration
	Tps        int

	MaxDelta        time.Duration
	MaxVirtualDelta time.Duration

	Scale  float64
	Paused bool

	MaxSteps int
}

func newClockSettings() ClockSettings {
	return ClockSettings{
		Mode:            ModeFixed,
		Now:             time.Now,
		FixedDelta:      DefaultFixedDelta,
		Tps:             DefaultTps,
		MaxDelta:        DefaultMaxDelta,
		MaxVirtualDelta: DefaultMaxVirtualDelta,
		Scale:           DefaultScale,
		Paused:          DefaultPaused,
		MaxSteps:        DefaultMaxSteps,
	}
}

// =======
// systems
// =======
func sanitizeClockSettings(w *ecs.World) {
	cs := ecs.EnsureResource(w, newClockSettings)

	// delta (master)
	if cs.FixedDelta <= 0 {
		// TODO: log that the user is dumb but now log depends on clock
		cs.FixedDelta = DefaultFixedDelta
	}

	derivedTps := int((time.Second + cs.FixedDelta/2) / cs.FixedDelta)
	cs.Tps = max(MinTps, min(derivedTps, MaxTps))

	// scale
	if math.IsNaN(cs.Scale) || math.IsInf(cs.Scale, 0) || cs.Scale <= 0 {
		cs.Scale = DefaultScale
	}

	// max delta
	if cs.MaxDelta <= 0 {
		cs.MaxDelta = DefaultMaxDelta
	}

	// virtual delta
	if cs.MaxVirtualDelta < 0 {
		cs.MaxVirtualDelta = DefaultMaxVirtualDelta
	}

	// max steps
	if cs.MaxSteps <= 0 {
		cs.MaxSteps = DefaultMaxSteps
	}
}

// ===
// API
// ===
func SetMode(w *ecs.World, m ClockMode) {
	cs := ecs.EnsureResource(w, newClockSettings)
	cs.Mode = m
}

func SetTPS(w *ecs.World, tps int) {
	if tps <= 0 {
		tps = DefaultTps
	}

	cs := ecs.EnsureResource(w, newClockSettings)
	cs.Mode = ModeFixed
	cs.FixedDelta = time.Second / time.Duration(tps)
	cs.Tps = tps
}

func SetFixedDelta(w *ecs.World, d time.Duration) {
	if d <= 0 {
		d = DefaultFixedDelta
	}

	cs := ecs.EnsureResource(w, newClockSettings)
	cs.Mode = ModeFixed
	cs.FixedDelta = d
	cs.Tps = int(math.Round(1.0 / d.Seconds()))
}

func SetScale(w *ecs.World, scale float64) {
	cs := ecs.EnsureResource(w, newClockSettings)
	cs.Scale = scale
}

func PauseClock(w *ecs.World, v bool) {
	cs := ecs.EnsureResource(w, newClockSettings)
	cs.Paused = v
}

func SetMaxSteps(w *ecs.World, n int) {
	cs := ecs.EnsureResource(w, newClockSettings)
	cs.MaxSteps = n
}

func SetMaxDelta(w *ecs.World, d time.Duration) {
	cs := ecs.EnsureResource(w, newClockSettings)
	cs.MaxDelta = d
}

func SetMaxVirtualDelta(w *ecs.World, d time.Duration) {
	cs := ecs.EnsureResource(w, newClockSettings)
	cs.MaxVirtualDelta = d
}
