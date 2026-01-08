package input

// ============
// mouse button
// ============
type MouseButton uint8

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

// =====
// mouse
// =====
type Mouse struct {
	Button  ButtonInput[MouseButton]
	CursorX AnalogInput
	CursorY AnalogInput
	Wheel   AnalogInput
}

func newMouse() Mouse {
	return Mouse{
		Button:  newButtonInput[MouseButton](),
		CursorX: newAnalogInput(AnalogAbsolute),
		CursorY: newAnalogInput(AnalogAbsolute),
		Wheel:   newAnalogInput(AnalogTransient),
	}
}

func (m *Mouse) begin() {
	m.Button.begin()
	m.CursorX.begin()
	m.CursorY.begin()
	m.Wheel.begin()
}

func (m *Mouse) compute() {
	m.Button.compute()
	m.CursorX.compute()
	m.CursorY.compute()
	m.Wheel.compute()
}

func (m *Mouse) end() {
	m.Button.end()
	m.CursorX.end()
	m.CursorY.end()
	m.Wheel.end()
}

func (m *Mouse) reset() {
	m.Button.reset()
	m.CursorX.reset()
	m.CursorY.reset()
	m.Wheel.reset()
}

// ======
// events
// ======
type MouseButtonEvent struct {
	Button  MouseButton
	Pressed bool
}

type MouseMoveEvent struct {
	X float64
	Y float64
}

type MouseWheelEvent struct {
	Delta float64
}
