package time

import (
	"time"

	"github.com/vistormu/xpeto/internal/core"
)

type Clock struct {
	accumulator float64
	lastTime    time.Time
}

func NewClock() *Clock {
	return &Clock{
		accumulator: 0,
		lastTime:    time.Now(),
	}
}

func (c *Clock) Update(ctx *core.Context) {
	t, ok := core.GetResource[*Time](ctx)
	if !ok {
		return
	}

	now := time.Now()
	if t.Frame == 0 {
		t.LastReal = now
	}
	rawDelta := now.Sub(t.LastReal)
	if t.maxDelta > 0 && rawDelta > t.maxDelta {
		rawDelta = t.maxDelta
	}
	t.Delta = rawDelta
	t.LastReal = now
	t.RealElapsed += rawDelta

	if t.Paused {
		t.ScaledDelta = 0
	} else {
		t.ScaledDelta = time.Duration(float64(rawDelta) * t.Scale)
		t.Elapsed += t.ScaledDelta
	}

	t.Frame++
}
