package ecs

import (
	"reflect"

	st "github.com/vistormu/xpeto/internal/structures"
)

type EntityManager struct {
	nextId     uint32
	storage    map[reflect.Type]map[Entity]any
	population *st.HashSet[Entity]
}

func NewEntityManager() *EntityManager {
	return &EntityManager{
		nextId:     0,
		storage:    make(map[reflect.Type]map[Entity]any),
		population: st.NewHashSet[Entity](),
	}
}

// ========
// entities
// ========
func (em *EntityManager) Create() Entity {
	entity := Entity{Id: em.nextId}
	em.nextId++

	em.population.Add(entity)

	return entity
}

func (em *EntityManager) Destroy(entity Entity) {
	if !em.population.Contains(entity) {
		return
	}

	em.population.Remove(entity)

	for _, storage := range em.storage {
		delete(storage, entity)
	}
}

func (em *EntityManager) DestroyAll() {
	for _, entity := range em.population.Values() {
		em.Destroy(entity)
	}
}

func (em *EntityManager) Query(f Filter) []Entity {
	var result []Entity
	for _, entity := range em.population.Values() {
		if f.Match(em, entity) {
			result = append(result, entity)
		}
	}
	return result
}

// ==========
// components
// ==========
func GetComponent[T any](em *EntityManager, id Entity) (T, bool) {
	var zero T
	component, ok := em.storage[reflect.TypeOf(zero)][id]
	if !ok {
		return zero, false
	}

	return component.(T), true
}

func AddComponent[T any](em *EntityManager, entity Entity, component T) {
	storage, ok := em.storage[reflect.TypeOf(component)]
	if !ok {
		storage = make(map[Entity]any)
		em.storage[reflect.TypeOf(component)] = storage
	}

	storage[entity] = component
}

func RemoveComponent[T any](em *EntityManager, entity Entity) {
	var zero T
	storage, ok := em.storage[reflect.TypeOf(zero)]
	if !ok {
		return
	}

	delete(storage, entity)
}
