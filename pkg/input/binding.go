package input

type BindingBuilder struct {
	action   Action
	bindings []Binding
}

func NewBinding(action Action) *BindingBuilder {
	return &BindingBuilder{
		action:   action,
		bindings: []Binding{},
	}
}

func (bb *BindingBuilder) Keyboard(key Key, trigger Trigger) *BindingBuilder {
	binding := Binding{
		Device:  KeyboardInput,
		Code:    InputCode{Key: key},
		Trigger: trigger,
	}
	bb.bindings = append(bb.bindings, binding)
	return bb
}

func (bb *BindingBuilder) Mouse(button MouseButton, trigger Trigger) *BindingBuilder {
	binding := Binding{
		Device:  MouseInput,
		Code:    InputCode{Mouse: button},
		Trigger: trigger,
	}
	bb.bindings = append(bb.bindings, binding)
	return bb
}

func (bb *BindingBuilder) Gamepad(id int, button GamepadButton, trigger Trigger) *BindingBuilder {
	binding := Binding{
		Device:  GamepadInput,
		Code:    InputCode{Gamepad: Gamepad{Id: id, Button: button}},
		Trigger: trigger,
	}
	bb.bindings = append(bb.bindings, binding)
	return bb
}

func (bb *BindingBuilder) GamepadAxis(id int, axis GamepadAxis, trigger Trigger) *BindingBuilder {
	binding := Binding{
		Device:  GamepadInput,
		Code:    InputCode{Gamepad: Gamepad{Id: id, Axis: axis}},
		Trigger: trigger,
	}
	bb.bindings = append(bb.bindings, binding)
	return bb
}

func (bb *BindingBuilder) Build(manager *Manager) {
	manager.RegisterBindings(bb.action, bb.bindings)
}
