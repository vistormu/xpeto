package scene

import (
	"reflect"

	"github.com/vistormu/xpeto/internal/errors"
	st "github.com/vistormu/xpeto/internal/structures"
)

type Manager struct {
	scenes map[reflect.Type]Scene
	stack  *st.UniqueStack[reflect.Type]
	active *st.HashSet[Scene]
}

func NewManager() *Manager {
	return &Manager{
		scenes: make(map[reflect.Type]Scene),
		stack:  st.NewUniqueStack[reflect.Type](),
		active: st.NewHashSet[Scene](),
	}
}

func (m *Manager) Register(scene Scene) {
	_, ok := m.scenes[reflect.TypeOf(scene)]
	if ok {
		return
	}

	m.scenes[reflect.TypeOf(scene)] = scene
}

func (m *Manager) Push(typ reflect.Type) {
	m.stack.Push(typ)
	scene, ok := m.scenes[typ]
	if !ok {
		errors.New(errors.SceneNotRegistered).With("type", typ).Print()
		return
	}
	m.active.Add(scene)
}

func (m *Manager) Pop() {
	if m.stack.IsEmpty() {
		errors.New(errors.SceneStackEmpty).Print()
		return
	}

	typ, _ := m.stack.Pop()
	scene, ok := m.scenes[typ]
	if !ok {
		errors.New(errors.SceneNotRegistered).With("type", typ).Print()
		return
	}
	m.active.Remove(scene)
}

func (m *Manager) Scene(typ reflect.Type) (Scene, bool) {
	scene, ok := m.scenes[typ]
	if !ok {
		return nil, false
	}
	return scene, true
}

func (m *Manager) Current() (Scene, bool) {
	if m.stack.IsEmpty() {
		return nil, false
	}

	typ, _ := m.stack.Peek()
	scene, ok := m.scenes[typ]
	return scene, ok
}

func (m *Manager) IsActive(scene Scene) bool {
	return m.active.Contains(scene)
}

func (m *Manager) Active() []Scene {
	return m.active.Values()
}

func (m *Manager) Clear() {
	m.stack.Clear()
	m.active.Clear()
}
