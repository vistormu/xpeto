package input

import (
	"testing"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/event"
	"github.com/vistormu/xpeto/core/schedule"
)

func TestButtonInput_FrameFlagsAndDuration(t *testing.T) {
	kb := newKeyboard()

	// frame 0 start
	kb.begin()
	kb.press(KeyA)

	if !kb.IsPressed(KeyA) {
		t.Fatalf("expected pressed")
	}
	if !kb.IsJustPressed(KeyA) {
		t.Fatalf("expected just pressed in same frame")
	}

	// frame 0 compute
	kb.compute()
	if got := kb.Duration(KeyA); got != 1 {
		t.Fatalf("expected duration 1 after compute, got %d", got)
	}

	// frame 1 start (flags must clear)
	kb.begin()
	if kb.IsJustPressed(KeyA) {
		t.Fatalf("expected just pressed cleared at next beginFrame")
	}
	if kb.IsJustReleased(KeyA) {
		t.Fatalf("expected just released cleared at next beginFrame")
	}

	// frame 1 compute
	kb.compute()
	if got := kb.Duration(KeyA); got != 2 {
		t.Fatalf("expected duration 2 after second compute, got %d", got)
	}

	// release in frame 1
	kb.release(KeyA)
	if kb.IsPressed(KeyA) {
		t.Fatalf("expected not pressed after release")
	}
	if !kb.IsJustReleased(KeyA) {
		t.Fatalf("expected just released in same frame")
	}
	if got := kb.Duration(KeyA); got != 0 {
		t.Fatalf("expected duration 0 after release, got %d", got)
	}

	// frame 2 start clears release flag
	kb.begin()
	if kb.IsJustReleased(KeyA) {
		t.Fatalf("expected just released cleared at next beginFrame")
	}
}

func TestAnalogInput_Transient_Behaviour(t *testing.T) {
	wheel := newAnalogInput(AnalogTransient)

	// frame start
	wheel.begin()

	// backend writes during frame
	wheel.add(2)
	wheel.add(3)
	if got := wheel.Value(); got != 5 {
		t.Fatalf("expected value 5 before compute, got %v", got)
	}

	// compute makes delta visible for the frame
	wheel.compute()
	if got := wheel.Delta(); got != 5 {
		t.Fatalf("expected delta 5 after compute, got %v", got)
	}

	// endFrame clears transient value
	wheel.end()
	if got := wheel.Value(); got != 0 {
		t.Fatalf("expected value 0 after endFrame, got %v", got)
	}
}

func TestAnalogInput_Absolute_DeltaAndPrevious(t *testing.T) {
	ax := newAnalogInput(AnalogAbsolute)

	// frame 0
	ax.begin()
	ax.set(10)
	ax.compute()
	ax.end()

	if got := ax.Delta(); got != 10 {
		t.Fatalf("expected delta 10, got %v", got)
	}
	if got := ax.Previous(); got != 10 {
		t.Fatalf("expected previous 10, got %v", got)
	}

	// frame 1
	ax.begin()
	ax.set(12)
	ax.compute()
	ax.end()

	if got := ax.Delta(); got != 2 {
		t.Fatalf("expected delta 2, got %v", got)
	}
	if got := ax.Previous(); got != 12 {
		t.Fatalf("expected previous 12, got %v", got)
	}
}

func TestWatchFocus_UnfocusedResetsAllInput(t *testing.T) {
	w := ecs.NewWorld()

	// event.GetEvents requires a running system resource.
	ecs.AddResource(w, schedule.RunningSystem{Id: 1, Label: "test"})

	// event plugin resource setup
	event.Pkg(w, schedule.NewScheduler())

	ecs.AddResource(w, newKeyboard())
	ecs.AddResource(w, newMouse())
	ecs.AddResource(w, newGamepads())

	kb, _ := ecs.GetResource[Keyboard](w)
	m, _ := ecs.GetResource[Mouse](w)
	gs, _ := ecs.GetResource[Gamepads](w)

	kb.press(KeyA)
	m.Button.press(MouseButtonLeft)
	m.Wheel.add(1)

	gs.connect(7, GamepadInfo{Name: "pad"})
	g, _ := gs.Get(7)
	g.Buttons.press(GamepadButtonStart)
	g.Axis(GamepadAxisLeftX).set(0.25)

	event.AddEvent(w, FocusChangedEvent{Focused: false})
	watchFocus(w)

	if kb.IsPressed(KeyA) || kb.IsJustPressed(KeyA) || kb.IsJustReleased(KeyA) {
		t.Fatalf("expected keyboard reset on focus lost")
	}
	if m.Button.IsPressed(MouseButtonLeft) || m.Button.IsJustPressed(MouseButtonLeft) || m.Button.IsJustReleased(MouseButtonLeft) {
		t.Fatalf("expected mouse buttons reset on focus lost")
	}
	if m.Wheel.Value() != 0 || m.Wheel.Delta() != 0 {
		t.Fatalf("expected wheel reset on focus lost, got value=%v delta=%v", m.Wheel.Value(), m.Wheel.Delta())
	}

	g2, ok := gs.Get(7)
	if !ok || g2 == nil {
		t.Fatalf("expected gamepad still present")
	}
	if g2.Buttons.IsPressed(GamepadButtonStart) || g2.Buttons.IsJustPressed(GamepadButtonStart) || g2.Buttons.IsJustReleased(GamepadButtonStart) {
		t.Fatalf("expected gamepad buttons reset on focus lost")
	}
	if got := g2.Axis(GamepadAxisLeftX).Delta(); got != 0 {
		t.Fatalf("expected axis delta reset on focus lost, got %v", got)
	}
}
