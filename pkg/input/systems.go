package input

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/event"
)

func beginFrame(w *ecs.World) {
	kb := ecs.EnsureResource(w, newKeyboard)
	kb.begin()

	m := ecs.EnsureResource(w, newMouse)
	m.begin()

	gs := ecs.EnsureResource(w, newGamepads)
	for g := range gs.Iter() {
		g.begin()
	}
}

func applyEvents(w *ecs.World) {
	// keyboard
	kb := ecs.EnsureResource(w, newKeyboard)
	if ev, ok := event.GetEvents[KeyEvent](w); ok {
		for _, e := range ev {
			if e.Pressed {
				kb.press(e.Key)
			} else {
				kb.release(e.Key)
			}
		}
	}

	// mouse
	m := ecs.EnsureResource(w, newMouse)
	if ev, ok := event.GetEvents[MouseButtonEvent](w); ok {
		for _, e := range ev {
			if e.Pressed {
				m.Button.press(e.Button)
			} else {
				m.Button.release(e.Button)
			}
		}
	}
	if ev, ok := event.GetEvents[MouseMoveEvent](w); ok && len(ev) != 0 {
		last := ev[len(ev)-1] // last wins
		m.CursorX.set(last.X)
		m.CursorY.set(last.Y)
	}
	if ev, ok := event.GetEvents[MouseWheelEvent](w); ok {
		for _, e := range ev {
			m.Wheel.add(e.Delta)
		}
	}

	// gamepads: connect/disconnect
	gs := ecs.EnsureResource(w, newGamepads)
	if ev, ok := event.GetEvents[GamepadConnectionEvent](w); ok {
		for _, e := range ev {
			switch e.Kind {
			case GamepadConnected:
				gs.connect(e.Id, e.Info)
			case GamepadDisconnected:
				gs.disconnect(e.Id)
			}
		}
	}

	// gamepads: buttons
	if ev, ok := event.GetEvents[GamepadButtonEvent](w); ok {
		for _, e := range ev {
			g := gs.Ensure(e.Id)
			if e.Pressed {
				g.Buttons.press(e.Button)
			} else {
				g.Buttons.release(e.Button)
			}
		}
	}

	// gamepads: axes
	if ev, ok := event.GetEvents[GamepadAxisEvent](w); ok {
		for _, e := range ev {
			g := gs.Ensure(e.Id)
			g.Axis(e.Axis).set(e.Value)
		}
	}
}

func compute(w *ecs.World) {
	kb := ecs.EnsureResource(w, newKeyboard)
	kb.compute()

	m := ecs.EnsureResource(w, newMouse)
	m.compute()

	gs := ecs.EnsureResource(w, newGamepads)
	for g := range gs.Iter() {
		g.compute()
	}
}

func endFrame(w *ecs.World) {
	kb := ecs.EnsureResource(w, newKeyboard)
	kb.end()

	m := ecs.EnsureResource(w, newMouse)
	m.end()

	gs := ecs.EnsureResource(w, newGamepads)
	for g := range gs.Iter() {
		g.end()
	}
}
