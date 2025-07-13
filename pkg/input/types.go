package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vistormu/xpeto/internal/core"
)

// ============
// input device
// ============
type InputDevice int

const (
	KeyboardInput InputDevice = iota
	MouseInput
	GamepadInput
	// TouchInput
)

// ======
// action
// ======
type Action = core.Handle[any]

// =======
// trigger
// =======
type Trigger int

const (
	Press Trigger = iota
	Release
	Maintain
	DoubleTap
)

// ==========
// input code
// ==========
type InputCode struct {
	Key     Key
	Mouse   MouseButton
	Gamepad Gamepad
}

// =======
// binding
// =======
type Binding struct {
	Device  InputDevice
	Code    InputCode
	Trigger Trigger
}

// ======
// inputs
// ======
// keys
type Key = ebiten.Key

var (
	// a to z
	KeyA = ebiten.KeyA
	KeyB = ebiten.KeyB
	KeyC = ebiten.KeyC
	KeyD = ebiten.KeyD
	KeyE = ebiten.KeyE
	KeyF = ebiten.KeyF
	KeyG = ebiten.KeyG
	KeyH = ebiten.KeyH
	KeyI = ebiten.KeyI
	KeyJ = ebiten.KeyJ
	KeyK = ebiten.KeyK
	KeyL = ebiten.KeyL
	KeyM = ebiten.KeyM
	KeyN = ebiten.KeyN
	KeyO = ebiten.KeyO
	KeyP = ebiten.KeyP
	KeyQ = ebiten.KeyQ
	KeyR = ebiten.KeyR
	KeyS = ebiten.KeyS
	KeyT = ebiten.KeyT
	KeyU = ebiten.KeyU
	KeyV = ebiten.KeyV
	KeyW = ebiten.KeyW
	KeyX = ebiten.KeyX
	KeyY = ebiten.KeyY
	KeyZ = ebiten.KeyZ

	// numbers
	Key0 = ebiten.Key0
	Key1 = ebiten.Key1
	Key2 = ebiten.Key2
	Key3 = ebiten.Key3
	Key4 = ebiten.Key4
	Key5 = ebiten.Key5
	Key6 = ebiten.Key6
	Key7 = ebiten.Key7
	Key8 = ebiten.Key8
	Key9 = ebiten.Key9

	// symbols
	KeyMinus        = ebiten.KeyMinus
	KeyEqual        = ebiten.KeyEqual
	KeyLeftBracket  = ebiten.KeyLeftBracket
	KeyRightBracket = ebiten.KeyRightBracket
	KeyBackslash    = ebiten.KeyBackslash
	KeySemicolon    = ebiten.KeySemicolon

	// special keys
	KeySpace     = ebiten.KeySpace
	KeyEnter     = ebiten.KeyEnter
	KeyEscape    = ebiten.KeyEscape
	KeyBackspace = ebiten.KeyBackspace
	KeyTab       = ebiten.KeyTab
	KeyShift     = ebiten.KeyShift
	KeyControl   = ebiten.KeyControl
	KeyAlt       = ebiten.KeyAlt
)

// =====
// mouse
// =====
type MouseButton = ebiten.MouseButton

var (
	MouseButtonLeft   = ebiten.MouseButtonLeft
	MouseButtonMiddle = ebiten.MouseButtonMiddle
	MouseButtonRight  = ebiten.MouseButtonRight
)

// =======
// gamepad
// =======
type GamepadButton = ebiten.GamepadButton
type GamepadAxis = ebiten.GamepadAxisType
type GamepadId = ebiten.GamepadID

type Gamepad struct {
	Id     int
	Button GamepadButton
	Axis   GamepadAxis
}

var (
	GamepadButtonA = ebiten.GamepadButton0
)
