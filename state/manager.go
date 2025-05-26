package state

import (
	st "github.com/vistormu/xpeto/internal/structures"
)

type Manager struct {
	nextId uint32
	active *st.HashSet[State]
}

func NewManager() *Manager {
	return &Manager{
		active: st.NewHashSet[State](),
	}
}

func (m *Manager) Register() State {
	state := State{Id: m.nextId, Version: 0}
	m.nextId++
	m.active.Add(state)

	return state
}

func (m *Manager) Add(state State) {
	m.active.Add(state)
}

func (m *Manager) Remove(state State) {
	m.active.Remove(state)
}

func (m *Manager) IsActive(state State) bool {
	return m.active.Contains(state)
}
