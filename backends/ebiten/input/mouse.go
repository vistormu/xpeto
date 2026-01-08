package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/event"
	"github.com/vistormu/xpeto/pkg/input"
)

// =====
// table
// =====
var mouseButtonMap = [...]input.MouseButton{
	ebiten.MouseButtonLeft:   input.MouseButtonLeft,
	ebiten.MouseButtonMiddle: input.MouseButtonMiddle,
	ebiten.MouseButtonRight:  input.MouseButtonRight,
}

func mapMouseButton(mb ebiten.MouseButton) (input.MouseButton, bool) {
	if mb < 0 || int(mb) >= len(mouseButtonMap) {
		return 0, false
	}
	return mouseButtonMap[int(mb)], true
}

// =====
// event
// =====
func emitMouse(w *ecs.World) {
	for b := ebiten.MouseButtonLeft; b <= ebiten.MouseButtonRight; b++ {
		mb, ok := mapMouseButton(b)
		if !ok {
			continue
		}
		if inpututil.IsMouseButtonJustPressed(b) {
			event.AddEvent(w, input.MouseButtonEvent{Button: mb, Pressed: true})
		}
		if inpututil.IsMouseButtonJustReleased(b) {
			event.AddEvent(w, input.MouseButtonEvent{Button: mb, Pressed: false})
		}
	}

	x, y := ebiten.CursorPosition()
	event.AddEvent(w, input.MouseMoveEvent{X: float64(x), Y: float64(y)})

	_, wy := ebiten.Wheel()
	if wy != 0 {
		event.AddEvent(w, input.MouseWheelEvent{Delta: wy})
	}
}
