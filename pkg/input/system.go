package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/event"
)

type actionState struct {
	pressed bool
}

type ActionSystem struct {
	states map[Action]*actionState
}

func NewActionSystem() *ActionSystem {
	return &ActionSystem{
		states: make(map[Action]*actionState),
	}
}

func (s *ActionSystem) Update(ctx *core.Context, dt float32) {
	am, _ := core.GetResource[*Manager](ctx)
	em, _ := core.GetResource[*event.Bus](ctx)

	for action, bindings := range am.Mappings() {
		state, ok := s.states[action]
		if !ok {
			state = &actionState{}
			s.states[action] = state
		}

		// check for active bindings
		anyPressed := false
		for _, binding := range bindings {
			var pressed bool
			switch binding.Device {
			case KeyboardInput:
				pressed = ebiten.IsKeyPressed(binding.Code.Key)

			case MouseInput:
				pressed = ebiten.IsMouseButtonPressed(binding.Code.Mouse)

			case GamepadInput:
				pressed = ebiten.IsGamepadButtonPressed(ebiten.GamepadID(binding.Code.Gamepad.Id), binding.Code.Gamepad.Button)

			}

			// edge detection
			justPressed := pressed && !state.pressed
			justReleased := !pressed && state.pressed
			maintained := pressed && state.pressed

			if justPressed && binding.Trigger == Press {
				em.Publish(ActionPress{Action: action})
			}
			if justReleased && binding.Trigger == Release {
				em.Publish(ActionRelease{Action: action})
			}
			if maintained && binding.Trigger == Maintain {
				em.Publish(ActionMaintain{Action: action})
			}

			if pressed {
				anyPressed = true
			}
		}

		// update state
		state.pressed = anyPressed
	}
}
