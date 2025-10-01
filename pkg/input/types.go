package input

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// ========
// keyboard
// ========
type Key = ebiten.Key
type Keyboard = ButtonInput[Key]

func newKeyboard() *Keyboard {
	return newButtonInput[Key]()
}

// =====
// mouse
// =====
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

// =======
// gamepad
// =======
type GamepadButton = ebiten.GamepadButton
type GamepadAxis = ebiten.GamepadAxisType
type Gamepad struct {
	Button ButtonInput[GamepadButton]
	Axis   AxisInput
}
