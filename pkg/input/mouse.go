package input

type MouseButton int

const (
	MouseButton0 MouseButton = iota
	MouseButton1
	MouseButton2
	MouseButton3
	MouseButton4
	MouseButtonLeft   = MouseButton0
	MouseButtonMiddle = MouseButton1
	MouseButtonRight  = MouseButton2
	MouseButtonMax    = MouseButton4
)

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
