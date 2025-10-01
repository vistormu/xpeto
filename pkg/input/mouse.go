package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/event"
)

func updateMouseInput(ctx *core.Context) {
	mouse := core.MustResource[*Mouse](ctx)
	eb := core.MustResource[*event.Bus](ctx)

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
			event.Publish(eb, MouseButtonJustPressed{Button: btn})
		}

		// falling edges
		if inpututil.IsMouseButtonJustReleased(btn) {
			mouse.Button.release(btn)
			event.Publish(eb, MouseButtonJustReleased{Button: btn})
		}
	}
}

func updateMouseCursor(ctx *core.Context) {
	mouse := core.MustResource[*Mouse](ctx)

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

func updateMouseWheel(ctx *core.Context) {
	mouse := core.MustResource[*Mouse](ctx)

	x, y := ebiten.Wheel()

	mouse.Wheel.X = x
	mouse.Wheel.Y = y
}
