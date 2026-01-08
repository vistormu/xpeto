package input

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/event"
)

// ======
// events
// ======
type FocusChangedEvent struct {
	Focused bool
}

func watchFocus(w *ecs.World) {
	ev, ok := event.GetEvents[FocusChangedEvent](w)
	if !ok || len(ev) == 0 {
		return
	}

	e := ev[len(ev)-1]

	if e.Focused {
		return
	}

	kb, ok := ecs.GetResource[Keyboard](w)
	if ok {
		kb.reset()
	}

	m, ok := ecs.GetResource[Mouse](w)
	if ok {
		m.reset()
	}

	gps, ok := ecs.GetResource[Gamepads](w)
	if ok {
		for gp := range gps.Iter() {
			gp.reset()
		}
	}
}
