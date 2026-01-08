package input

import (
	"slices"
)

// =======
// buttons
// =======
type GamepadButton uint16

const (
	GamepadButtonSouth GamepadButton = iota
	GamepadButtonEast
	GamepadButtonWest
	GamepadButtonNorth
	GamepadButtonL1
	GamepadButtonR1
	GamepadButtonL2
	GamepadButtonR2
	GamepadButtonSelect
	GamepadButtonStart
	GamepadButtonLStick
	GamepadButtonRStick
	GamepadButtonDpadUp
	GamepadButtonDpadDown
	GamepadButtonDpadLeft
	GamepadButtonDpadRight
)

// ====
// axes
// ====
type GamepadAxis uint16

const (
	GamepadAxisLeftX GamepadAxis = iota
	GamepadAxisLeftY
	GamepadAxisRightX
	GamepadAxisRightY
	GamepadAxisTriggerLeft
	GamepadAxisTriggerRight
)

// ========
// metadata
// ========
type GamepadId uint32

type GamepadInfo struct {
	Name      string
	VendorId  uint16
	ProductId uint16
	Mapping   string
}

// ======
// events
// ======
type GamepadEventKind uint8

const (
	GamepadConnected GamepadEventKind = iota
	GamepadDisconnected
)

type GamepadConnectionEvent struct {
	Id   GamepadId
	Kind GamepadEventKind
	Info GamepadInfo
}

type GamepadButtonEvent struct {
	Id      GamepadId
	Button  GamepadButton
	Pressed bool
}

type GamepadAxisEvent struct {
	Id    GamepadId
	Axis  GamepadAxis
	Value float64
}

// =======
// gamepad
// =======
type Gamepad struct {
	Buttons ButtonInput[GamepadButton]
	Axes    map[GamepadAxis]*AnalogInput
}

func newGamepad() Gamepad {
	return Gamepad{
		Buttons: newButtonInput[GamepadButton](),
		Axes:    make(map[GamepadAxis]*AnalogInput),
	}
}

func (g *Gamepad) begin() {
	g.Buttons.begin()
	for _, ai := range g.Axes {
		ai.begin()
	}
}

func (g *Gamepad) compute() {
	g.Buttons.compute()
	for _, ai := range g.Axes {
		ai.compute()
	}
}

func (g *Gamepad) end() {
	g.Buttons.end()
	for _, ai := range g.Axes {
		ai.end()
	}
}

func (g *Gamepad) reset() {
	g.Buttons.reset()
	for _, ai := range g.Axes {
		ai.reset()
	}
}

func (g *Gamepad) Axis(id GamepadAxis) *AnalogInput {
	a, ok := g.Axes[id]
	if ok && a != nil {
		return a
	}

	newA := newAnalogInput(AnalogAbsolute)
	g.Axes[id] = &newA

	return &newA
}

func (g *Gamepad) AxisTransient(id GamepadAxis) *AnalogInput {
	a, ok := g.Axes[id]
	if ok && a != nil {
		return a
	}

	newA := newAnalogInput(AnalogTransient)
	g.Axes[id] = &newA

	return &newA
}

// ========
// resource
// ========
type Gamepads struct {
	pads map[GamepadId]*Gamepad
	info map[GamepadId]GamepadInfo
}

func newGamepads() Gamepads {
	return Gamepads{
		pads: make(map[GamepadId]*Gamepad),
		info: make(map[GamepadId]GamepadInfo),
	}
}

func (gs *Gamepads) connect(id GamepadId, info GamepadInfo) bool {
	_, existed := gs.pads[id]
	gs.Ensure(id)
	gs.info[id] = info

	return !existed
}

func (gs *Gamepads) disconnect(id GamepadId) bool {
	_, ok := gs.pads[id]
	if !ok {
		return false
	}

	delete(gs.pads, id)
	delete(gs.info, id)

	return true
}

// ===
// API
// ===
func (gs *Gamepads) Ids() []GamepadId {
	ids := make([]GamepadId, 0, len(gs.pads))

	for id := range gs.pads {
		ids = append(ids, id)
	}

	slices.Sort(ids)

	return ids
}

func (gs *Gamepads) Iter() func(func(*Gamepad) bool) {
	return func(yield func(*Gamepad) bool) {
		for _, g := range gs.pads {
			if g == nil {
				continue
			}
			if !yield(g) {
				return
			}
		}
	}
}

func (gs *Gamepads) Has(id GamepadId) bool {
	_, ok := gs.pads[id]
	return ok
}

func (gs *Gamepads) Get(id GamepadId) (*Gamepad, bool) {
	g, ok := gs.pads[id]
	return g, ok
}

func (gs *Gamepads) Info(id GamepadId) (GamepadInfo, bool) {
	inf, ok := gs.info[id]
	return inf, ok
}

func (gs *Gamepads) Ensure(id GamepadId) *Gamepad {
	if g, ok := gs.pads[id]; ok && g != nil {
		return g
	}

	g := newGamepad()
	gs.pads[id] = &g

	return gs.pads[id]
}
