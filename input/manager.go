package input

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Manager struct {
	keys []Key
}

func NewManager() *Manager {
	return &Manager{
		keys: make([]Key, 0),
	}
}

func (m *Manager) Register(keys ...Key) {
	m.keys = append(m.keys, keys...)
}

func (m *Manager) Keys() []Key {
	return m.keys
}

func (m *Manager) IsPressed(key Key) bool {
	return ebiten.IsKeyPressed(key)
}
