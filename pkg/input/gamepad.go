package input

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type GamepadButton = ebiten.GamepadButton

type GamepadAxis = ebiten.GamepadAxisType

type Gamepad struct {
	Button ButtonInput[GamepadButton]
	Axis   AxisInput
}
