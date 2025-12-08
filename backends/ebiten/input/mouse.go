package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/pkg/input"
)

var mouseButtonTable [5]input.MouseButton

func init() {
	for i := range mouseButtonTable {
		mouseButtonTable[i] = input.MouseButton0
	}

	mouseButtonTable[ebiten.MouseButton0] = input.MouseButton0
	mouseButtonTable[ebiten.MouseButton1] = input.MouseButton1
	mouseButtonTable[ebiten.MouseButton2] = input.MouseButton2
	mouseButtonTable[ebiten.MouseButton3] = input.MouseButton3
	mouseButtonTable[ebiten.MouseButton4] = input.MouseButton4
}

func updateMouseInput(w *ecs.World) {
	mouse, _ := ecs.GetResource[input.Mouse](w)

	mouse.Button.Clear()

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
			mouse.Button.Press(mouseButtonTable[btn])
		}

		// falling edges
		if inpututil.IsMouseButtonJustReleased(btn) {
			mouse.Button.Release(mouseButtonTable[btn])
		}
	}
}

func updateMouseCursor(w *ecs.World) {
	mouse, _ := ecs.GetResource[input.Mouse](w)

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
	mouse, _ := ecs.GetResource[input.Mouse](w)

	x, y := ebiten.Wheel()

	mouse.Wheel.X = x
	mouse.Wheel.Y = y
}
