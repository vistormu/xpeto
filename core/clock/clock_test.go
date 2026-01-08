package clock

import (
	"testing"
	"time"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

type seqNow struct {
	times []time.Time
	i     int
}

func (s *seqNow) Now() time.Time {
	if len(s.times) == 0 {
		return time.Unix(0, 0)
	}
	if s.i >= len(s.times) {
		return s.times[len(s.times)-1]
	}
	t := s.times[s.i]
	s.i++
	return t
}

func newTestWorld(t *testing.T, times ...time.Time) *ecs.World {
	t.Helper()

	w := ecs.NewWorld()
	sch := schedule.NewScheduler()
	Pkg(w, sch)

	// deterministic Now
	cs, _ := ecs.GetResource[ClockSettings](w)
	seq := &seqNow{times: times}
	cs.Now = seq.Now

	// baseline runtime clocks
	base := time.Unix(0, 0)
	if len(times) > 0 {
		base = times[0]
	}

	rc, _ := ecs.GetResource[RealClock](w)
	rc.Start = base
	rc.Last = base
	rc.Now = base
	rc.Delta = 0
	rc.Clamped = 0
	rc.Elapsed = 0
	rc.Frame = 0

	vc, _ := ecs.GetResource[VirtualClock](w)
	vc.Start = base
	vc.Delta = 0
	vc.Elapsed = 0
	vc.Frame = 0

	fc, _ := ecs.GetResource[FixedClock](w)
	fc.Accumulator = 0
	fc.Steps = 0
	fc.Frame = 0
	fc.Alpha = 0

	return w
}

// step calls tick once.
func step(w *ecs.World) { tick(w) }

func TestTick_NegativeTimeJumpClampsToZero(t *testing.T) {
	t0 := time.Unix(10, 0)
	tBack := time.Unix(9, 0)

	w := newTestWorld(t, t0, tBack)

	step(w) // now=t0, delta 0
	step(w) // now=tBack, delta clamped to 0

	rc, _ := ecs.GetResource[RealClock](w)
	if rc.Delta != 0 {
		t.Fatalf("expected real delta 0 on backward time jump, got %v", rc.Delta)
	}
}

func TestTick_RealClampScaleAndVirtualClamp(t *testing.T) {
	t0 := time.Unix(0, 0)
	t1 := t0.Add(50 * time.Millisecond)

	w := newTestWorld(t, t0, t1)

	cs, _ := ecs.GetResource[ClockSettings](w)
	cs.MaxDelta = 10 * time.Millisecond
	cs.Scale = 2
	cs.MaxVirtualDelta = 15 * time.Millisecond

	step(w) // warm up, delta 0
	step(w) // delta 50ms

	rc, _ := ecs.GetResource[RealClock](w)
	if rc.Clamped != 10*time.Millisecond {
		t.Fatalf("expected clamped 10ms, got %v", rc.Clamped)
	}

	vc, _ := ecs.GetResource[VirtualClock](w)
	// 10ms * 2 = 20ms -> clamped to 15ms
	if vc.Delta != 15*time.Millisecond {
		t.Fatalf("expected virtual delta 15ms, got %v", vc.Delta)
	}
}

func TestTick_PausedOrZeroScaleFreezesVirtual(t *testing.T) {
	t0 := time.Unix(0, 0)
	t1 := t0.Add(16 * time.Millisecond)

	w := newTestWorld(t, t0, t1)

	cs, _ := ecs.GetResource[ClockSettings](w)
	cs.Paused = true
	cs.Scale = 10

	step(w)
	step(w)

	vc, _ := ecs.GetResource[VirtualClock](w)
	if vc.Delta != 0 || vc.Elapsed != 0 {
		t.Fatalf("expected virtual clock frozen, got delta=%v elapsed=%v", vc.Delta, vc.Elapsed)
	}
}

func TestTick_FixedStepsAndAlpha(t *testing.T) {
	t0 := time.Unix(0, 0)
	t1 := t0.Add(40 * time.Millisecond)

	w := newTestWorld(t, t0, t1)

	cs, _ := ecs.GetResource[ClockSettings](w)
	cs.Mode = ModeFixed
	cs.FixedDelta = 16 * time.Millisecond
	cs.MaxSteps = 8
	cs.Scale = 1
	cs.MaxDelta = 0

	step(w) // warm up
	step(w) // +40ms

	fc, _ := ecs.GetResource[FixedClock](w)
	if fc.Delta != 16*time.Millisecond {
		t.Fatalf("expected fixed delta 16ms, got %v", fc.Delta)
	}
	if fc.Steps != 2 {
		t.Fatalf("expected 2 fixed steps (40/16), got %d", fc.Steps)
	}
	// remainder 8ms -> alpha 0.5
	if fc.Alpha < 0.49 || fc.Alpha > 0.51 {
		t.Fatalf("expected alpha around 0.5, got %v", fc.Alpha)
	}
	if fc.Accumulator != 8*time.Millisecond {
		t.Fatalf("expected accumulator 8ms, got %v", fc.Accumulator)
	}
}

func TestTick_MaxStepsLimitsWork(t *testing.T) {
	t0 := time.Unix(0, 0)
	t1 := t0.Add(200 * time.Millisecond)

	w := newTestWorld(t, t0, t1)

	cs, _ := ecs.GetResource[ClockSettings](w)
	cs.Mode = ModeFixed
	cs.FixedDelta = 16 * time.Millisecond
	cs.MaxSteps = 1
	cs.Scale = 1

	step(w)
	step(w)

	fc, _ := ecs.GetResource[FixedClock](w)
	if fc.Steps != 1 {
		t.Fatalf("expected steps capped to 1, got %d", fc.Steps)
	}
	if fc.Accumulator <= 0 {
		t.Fatalf("expected leftover accumulator after cap, got %v", fc.Accumulator)
	}
	if fc.Alpha < 0 || fc.Alpha >= 1 {
		t.Fatalf("expected alpha in [0,1), got %v", fc.Alpha)
	}
}

func TestConditions_EveryNFrames(t *testing.T) {
	t0 := time.Unix(0, 0)
	t1 := t0.Add(1 * time.Millisecond)
	t2 := t1.Add(1 * time.Millisecond)

	w := newTestWorld(t, t0, t1, t2)

	cond := EveryNFrames(2)

	step(w) // frame 1
	if cond(w) {
		t.Fatalf("expected EveryNFrames(2) false on frame 1")
	}

	step(w) // frame 2
	if !cond(w) {
		t.Fatalf("expected EveryNFrames(2) true on frame 2")
	}
}

func TestConditions_OnceAfterElapsed(t *testing.T) {
	t0 := time.Unix(0, 0)
	t1 := t0.Add(10 * time.Millisecond)
	t2 := t1.Add(10 * time.Millisecond)

	w := newTestWorld(t, t0, t1, t2)

	cs, _ := ecs.GetResource[ClockSettings](w)
	cs.Scale = 1
	cs.Paused = false

	cond := OnceAfterElapsed(15 * time.Millisecond)

	step(w) // warm up (elapsed 0)
	step(w) // +10ms
	if cond(w) {
		t.Fatalf("expected false before elapsed reaches 15ms")
	}

	step(w) // +10ms => elapsed 20ms
	if !cond(w) {
		t.Fatalf("expected true once after elapsed reaches 15ms")
	}
	if cond(w) {
		t.Fatalf("expected false after once has fired")
	}
}

func TestConditions_EveryDurationPhase(t *testing.T) {
	t0 := time.Unix(0, 0)
	t1 := t0.Add(10 * time.Millisecond)
	t2 := t1.Add(10 * time.Millisecond)
	t3 := t2.Add(35 * time.Millisecond)

	w := newTestWorld(t, t0, t1, t2, t3)

	cs, _ := ecs.GetResource[ClockSettings](w)
	cs.Scale = 1
	cs.Paused = false

	cond := EveryDuration(20 * time.Millisecond)

	step(w) // warm up (elapsed 0)

	step(w) // +10ms => elapsed 10ms
	if cond(w) {
		t.Fatalf("expected false before 20ms")
	}

	step(w) // +10ms => elapsed 20ms
	if !cond(w) {
		t.Fatalf("expected true at 20ms")
	}

	step(w) // +35ms => elapsed 55ms, should fire (crossed 40ms)
	if !cond(w) {
		t.Fatalf("expected true on long frame crossing one or more periods")
	}
}
