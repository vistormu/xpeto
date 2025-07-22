package ecs

import (
	"reflect"

	"github.com/vistormu/xpeto/internal/core"
)

type Manager struct {
	nextId     uint32
	storage    map[reflect.Type]map[Entity]any
	population *core.HashSet[Entity]
}

func NewManager() *Manager {
	return &Manager{
		nextId:     1,
		storage:    make(map[reflect.Type]map[Entity]any),
		population: core.NewHashSet[Entity](),
	}
}

// ========
// entities
// ========
func (m *Manager) Create(archetype Archetype) Entity {
	entity := Entity{Id: m.nextId}
	m.nextId++

	m.population.Add(entity)

	if archetype != nil {
		for _, component := range archetype.Components() {
			AddComponent(m, entity, component)
		}
	}

	return entity
}

func (m *Manager) Destroy(entity Entity) {
	if !m.population.Contains(entity) {
		return
	}

	m.population.Remove(entity)

	for _, storage := range m.storage {
		delete(storage, entity)
	}
}

func (m *Manager) DestroyAll() {
	for _, entity := range m.population.Values() {
		m.Destroy(entity)
	}
}

func (m *Manager) Query(f Filter) []Entity {
	var result []Entity
	for _, entity := range m.population.Values() {
		if f.Match(m, entity) {
			result = append(result, entity)
		}
	}
	return result
}

// ==========
// components
// ==========
func GetComponent[T any](m *Manager, entity Entity) (T, bool) {
	var zero T
	component, ok := m.storage[reflect.TypeOf(zero)][entity]
	if !ok {
		return zero, false
	}

	return component.(T), true
}

func AddComponent[T any](m *Manager, entity Entity, component T) {
	storage, ok := m.storage[reflect.TypeOf(component)]
	if !ok {
		storage = make(map[Entity]any)
		m.storage[reflect.TypeOf(component)] = storage
	}

	storage[entity] = component
}

func RemoveComponent[T any](m *Manager, entity Entity) {
	var zero T
	storage, ok := m.storage[reflect.TypeOf(zero)]
	if !ok {
		return
	}

	delete(storage, entity)
}
