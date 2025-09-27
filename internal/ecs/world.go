package ecs

import (
	"reflect"

	"github.com/vistormu/xpeto/internal/core"
)

type World struct {
	nextId     int
	storage    map[reflect.Type]map[Entity]any
	population *core.HashSet[Entity]
}

func NewWorld() *World {
	return &World{
		nextId:     1,
		storage:    make(map[reflect.Type]map[Entity]any),
		population: core.NewHashSet[Entity](),
	}
}

// ========
// entities
// ========
func (w *World) Create() Entity {
	entity := Entity{Number: w.nextId}
	w.nextId++

	w.population.Add(entity)

	return entity
}

func (w *World) Destroy(entity Entity) {
	if !w.population.Contains(entity) {
		return
	}

	w.population.Remove(entity)

	for _, storage := range w.storage {
		delete(storage, entity)
	}
}

func (w *World) DestroyAll() {
	for _, entity := range w.population.Values() {
		w.Destroy(entity)
	}
}

func (w *World) Query(f Filter) []Entity {
	var result []Entity
	for _, entity := range w.population.Values() {
		if f.Match(w, entity) {
			result = append(result, entity)
		}
	}
	return result
}

// ==========
// components
// ==========
func GetComponent[T any](w *World, entity Entity) (T, bool) {
	var zero T
	component, ok := w.storage[reflect.TypeOf(zero)][entity]
	if !ok {
		return zero, false
	}

	return component.(T), true
}

func AddComponent[T any](w *World, entity Entity, component T) {
	storage, ok := w.storage[reflect.TypeOf(component)]
	if !ok {
		storage = make(map[Entity]any)
		w.storage[reflect.TypeOf(component)] = storage
	}

	storage[entity] = component
}

func RemoveComponent[T any](w *World, entity Entity) {
	var zero T
	storage, ok := w.storage[reflect.TypeOf(zero)]
	if !ok {
		return
	}

	delete(storage, entity)
}
