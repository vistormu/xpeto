package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/vistormu/xpeto/core/ecs"
	// "github.com/vistormu/xpeto/internal/event"
)

type MouseButton = ebiten.MouseButton
type Mouse struct {
	Button *ButtonInput[MouseButton]
	Cursor *CursorInput
	Wheel  *WheelInput
}

func newMouse() *Mouse {
	return &Mouse{
		Button: newButtonInput[MouseButton](),
		Cursor: &CursorInput{},
		Wheel:  &WheelInput{},
	}
}

func updateMouseInput(w *ecs.World) {
	mouse, _ := ecs.GetResource[Mouse](w)

	mouse.Button.clear()

	for _, btn := range []ebiten.MouseButton{
		ebiten.MouseButton0,
		ebiten.MouseButton1,
		ebiten.MouseButton2,
		ebiten.MouseButton3,
		ebiten.MouseButton4,
		ebiten.MouseButtonLeft,
		ebiten.MouseButtonMax,
		ebiten.MouseButtonMiddle,
		ebiten.MouseButtonRight,
	} {
		// rising edges
		if inpututil.IsMouseButtonJustPressed(btn) {
			mouse.Button.press(btn)
			// event.Publish(eb, MouseButtonJustPressed{Button: btn})
		}

		// falling edges
		if inpututil.IsMouseButtonJustReleased(btn) {
			mouse.Button.release(btn)
			// event.Publish(eb, MouseButtonJustReleased{Button: btn})
		}
	}
}

func updateMouseCursor(w *ecs.World) {
	mouse, _ := ecs.GetResource[Mouse](w)

	x, y := ebiten.CursorPosition()

	prevX := mouse.Cursor.X
	prevY := mouse.Cursor.Y

	mouse.Cursor.X = x
	mouse.Cursor.Y = y
	mouse.Cursor.PrevX = prevX
	mouse.Cursor.PrevY = prevY
	mouse.Cursor.Dx = x - prevX
	mouse.Cursor.Dy = y - prevY
}

func updateMouseWheel(w *ecs.World) {
	mouse, _ := ecs.GetResource[Mouse](w)

	x, y := ebiten.Wheel()

	mouse.Wheel.X = x
	mouse.Wheel.Y = y
}
