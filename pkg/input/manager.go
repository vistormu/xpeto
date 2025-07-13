package input

type Manager struct {
	nextId   uint32
	mappings map[Action][]Binding
}

func NewManager() *Manager {
	return &Manager{
		nextId:   1,
		mappings: make(map[Action][]Binding),
	}
}

func (m *Manager) Register() Action {
	action := Action{Id: m.nextId, Version: 0}
	m.nextId++

	return action
}

func (m *Manager) RegisterBindings(action Action, bindings []Binding) {
	_, ok := m.mappings[action]
	if !ok {
		m.mappings[action] = []Binding{}
	}

	m.mappings[action] = append(m.mappings[action], bindings...)
}

func (m *Manager) Mappings() map[Action][]Binding {
	return m.mappings
}
