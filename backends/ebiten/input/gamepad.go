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
var stdButtonMap = func() []uint16 {
	m := make([]uint16, int(ebiten.StandardGamepadButtonMax)+1)
	set := func(b ebiten.StandardGamepadButton, v input.GamepadButton) { m[int(b)] = uint16(v) + 1 }

	set(ebiten.StandardGamepadButtonRightBottom, input.GamepadButtonSouth)
	set(ebiten.StandardGamepadButtonRightRight, input.GamepadButtonEast)
	set(ebiten.StandardGamepadButtonRightLeft, input.GamepadButtonWest)
	set(ebiten.StandardGamepadButtonRightTop, input.GamepadButtonNorth)

	set(ebiten.StandardGamepadButtonFrontTopLeft, input.GamepadButtonL1)
	set(ebiten.StandardGamepadButtonFrontTopRight, input.GamepadButtonR1)
	set(ebiten.StandardGamepadButtonFrontBottomLeft, input.GamepadButtonL2)
	set(ebiten.StandardGamepadButtonFrontBottomRight, input.GamepadButtonR2)

	set(ebiten.StandardGamepadButtonCenterLeft, input.GamepadButtonSelect)
	set(ebiten.StandardGamepadButtonCenterRight, input.GamepadButtonStart)

	set(ebiten.StandardGamepadButtonLeftStick, input.GamepadButtonLStick)
	set(ebiten.StandardGamepadButtonRightStick, input.GamepadButtonRStick)

	set(ebiten.StandardGamepadButtonLeftTop, input.GamepadButtonDpadUp)
	set(ebiten.StandardGamepadButtonLeftBottom, input.GamepadButtonDpadDown)
	set(ebiten.StandardGamepadButtonLeftLeft, input.GamepadButtonDpadLeft)
	set(ebiten.StandardGamepadButtonLeftRight, input.GamepadButtonDpadRight)

	return m
}()

func mapStdButton(b ebiten.StandardGamepadButton) (input.GamepadButton, bool) {
	if b < 0 || int(b) >= len(stdButtonMap) {
		return 0, false
	}
	v := stdButtonMap[int(b)]
	if v == 0 {
		return 0, false
	}
	return input.GamepadButton(v - 1), true
}

func mapRawCommon(i int) (input.GamepadButton, bool) {
	// common raw layout
	switch i {
	case 0:
		return input.GamepadButtonSouth, true
	case 1:
		return input.GamepadButtonEast, true
	case 2:
		return input.GamepadButtonWest, true
	case 3:
		return input.GamepadButtonNorth, true
	case 4:
		return input.GamepadButtonL1, true
	case 5:
		return input.GamepadButtonR1, true
	case 6:
		return input.GamepadButtonSelect, true
	case 7:
		return input.GamepadButtonStart, true
	case 8:
		return input.GamepadButtonLStick, true
	case 9:
		return input.GamepadButtonRStick, true
	default:
		return 0, false
	}
}

// ======
// buffer
// ======
type gamepadState struct {
	connectedGamepads     []ebiten.GamepadID
	justConnectedGamepads []ebiten.GamepadID

	justPressedStdButtons  []ebiten.StandardGamepadButton
	justReleasedStdButtons []ebiten.StandardGamepadButton

	justPressedRawButtons  []ebiten.GamepadButton
	justReleasedRawButtons []ebiten.GamepadButton
}

func newGamepadState() gamepadState {
	return gamepadState{
		justConnectedGamepads:  make([]ebiten.GamepadID, 0),
		justPressedStdButtons:  make([]ebiten.StandardGamepadButton, 0),
		justReleasedStdButtons: make([]ebiten.StandardGamepadButton, 0),
		justPressedRawButtons:  make([]ebiten.GamepadButton, 0),
		justReleasedRawButtons: make([]ebiten.GamepadButton, 0),
	}
}

// =====
// event
// =====
func emitGamepads(w *ecs.World) {
	st := ecs.EnsureResource(w, newGamepadState)

	// connect
	st.justConnectedGamepads = inpututil.AppendJustConnectedGamepadIDs(st.justConnectedGamepads[:0])
	for _, id := range st.justConnectedGamepads {
		event.AddEvent(w, input.GamepadConnectionEvent{
			Id:   input.GamepadId(uint32(id)),
			Kind: input.GamepadConnected,
			Info: input.GamepadInfo{Name: ebiten.GamepadName(id)},
		})
	}

	// disconnect + state
	st.connectedGamepads = ebiten.AppendGamepadIDs(st.connectedGamepads[:0])
	for _, id := range st.connectedGamepads {
		gid := input.GamepadId(uint32(id))

		if inpututil.IsGamepadJustDisconnected(id) {
			event.AddEvent(w, input.GamepadConnectionEvent{Id: gid, Kind: input.GamepadDisconnected})
			continue
		}

		if ebiten.IsStandardGamepadLayoutAvailable(id) {
			// standard buttons
			st.justPressedStdButtons = inpututil.AppendJustPressedStandardGamepadButtons(id, st.justPressedStdButtons[:0])
			for _, b := range st.justPressedStdButtons {
				if bb, ok := mapStdButton(b); ok {
					event.AddEvent(w, input.GamepadButtonEvent{Id: gid, Button: bb, Pressed: true})
				}
			}

			st.justReleasedStdButtons = inpututil.AppendJustReleasedStandardGamepadButtons(id, st.justReleasedStdButtons[:0])
			for _, b := range st.justReleasedStdButtons {
				if bb, ok := mapStdButton(b); ok {
					event.AddEvent(w, input.GamepadButtonEvent{Id: gid, Button: bb, Pressed: false})
				}
			}

			// standard axes (absolute)
			event.AddEvent(w, input.GamepadAxisEvent{Id: gid, Axis: input.GamepadAxisLeftX, Value: ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal)})
			event.AddEvent(w, input.GamepadAxisEvent{Id: gid, Axis: input.GamepadAxisLeftY, Value: ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical)})
			event.AddEvent(w, input.GamepadAxisEvent{Id: gid, Axis: input.GamepadAxisRightX, Value: ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisRightStickHorizontal)})
			event.AddEvent(w, input.GamepadAxisEvent{Id: gid, Axis: input.GamepadAxisRightY, Value: ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisRightStickVertical)})
			// event.AddEvent(w, input.GamepadAxisEvent{Id: gid, Axis: input.GamepadAxisTriggerLeft, Value: ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftTrigger)})
			// event.AddEvent(w, input.GamepadAxisEvent{Id: gid, Axis: input.GamepadAxisTriggerRight, Value: ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisRightTrigger)})
			continue
		}

		st.justPressedRawButtons = inpututil.AppendJustPressedGamepadButtons(id, st.justPressedRawButtons[:0])
		for _, b := range st.justPressedRawButtons {
			if bb, ok := mapRawCommon(int(b)); ok {
				event.AddEvent(w, input.GamepadButtonEvent{Id: gid, Button: bb, Pressed: true})
			}
		}

		st.justReleasedRawButtons = inpututil.AppendJustReleasedGamepadButtons(id, st.justReleasedRawButtons[:0])
		for _, b := range st.justReleasedRawButtons {
			if bb, ok := mapRawCommon(int(b)); ok {
				event.AddEvent(w, input.GamepadButtonEvent{Id: gid, Button: bb, Pressed: false})
			}
		}

		n := ebiten.GamepadAxisCount(id)
		get := func(i int) float64 {
			if i < 0 || i >= n {
				return 0
			}
			return ebiten.GamepadAxisValue(id, i)
		}

		event.AddEvent(w, input.GamepadAxisEvent{Id: gid, Axis: input.GamepadAxisLeftX, Value: get(0)})
		event.AddEvent(w, input.GamepadAxisEvent{Id: gid, Axis: input.GamepadAxisLeftY, Value: get(1)})
		event.AddEvent(w, input.GamepadAxisEvent{Id: gid, Axis: input.GamepadAxisRightX, Value: get(2)})
		event.AddEvent(w, input.GamepadAxisEvent{Id: gid, Axis: input.GamepadAxisRightY, Value: get(3)})
		event.AddEvent(w, input.GamepadAxisEvent{Id: gid, Axis: input.GamepadAxisTriggerLeft, Value: get(4)})
		event.AddEvent(w, input.GamepadAxisEvent{Id: gid, Axis: input.GamepadAxisTriggerRight, Value: get(5)})
	}
}
