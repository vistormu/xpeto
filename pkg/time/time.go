package timer

import (
	"time"
)

type System struct {
	accumulator float64
	lastTime    time.Time
	fixedDelta  float64
}

func NewSystem() *System {
	return &System{
		accumulator: 0,
		lastTime:    time.Now(),
		fixedDelta:  1.0 / 60.0,
	}
}

func (s *System) Update() {
	now := time.Now()
	frameTime := now.Sub(s.lastTime).Seconds()
	s.lastTime = now

	if frameTime > 0.25 {
		frameTime = 0.25
	}
	s.accumulator += frameTime

	steps := int(s.accumulator / s.fixedDelta)
	if steps > 0 {
		s.accumulator -= float64(steps) * s.fixedDelta
	}

}
