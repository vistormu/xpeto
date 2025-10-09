package ecs

import (
	"reflect"

	"github.com/vistormu/go-dsa/hashmap"
)

type World struct {
	population *population
	registry   *hashmap.TypeMap
	resources  *hashmap.TypeMap
}

func NewWorld() *World {
	w := &World{
		population: newPopulation(),
		registry:   hashmap.NewTypeMap(),
		resources:  hashmap.NewTypeMap(),
	}

	AddResource(w, &systemId{})

	return w
}

// ===
// API
// ===
// entities
func AddEntity(w *World) Entity {
	return w.population.add()
}

func RemoveEntity(w *World, e Entity) bool {
	ok := w.population.remove(e)
	if !ok {
		return false
	}

	removeComponents(w.registry, e)

	return true
}

func HasEntity(w *World, e Entity) bool {
	return w.population.has(e)
}

// components
func AddComponent[T any](w *World, e Entity, c T) bool {
	ok := w.population.has(e)
	if !ok {
		return false
	}

	getStore[T](w.registry).add(e, c)

	return true
}

func GetComponent[T any](w *World, e Entity) (*T, bool) {
	ok := w.population.has(e)
	if !ok {
		return nil, false
	}

	return getStore[T](w.registry).get(e)
}

func RemoveComponent[T any](w *World, e Entity) bool {
	ok := w.population.has(e)
	if !ok {
		return false
	}

	return getStore[T](w.registry).remove(e)
}

// resources
func AddResource[T any](w *World, r T) {
	hashmap.Add(w.resources, r)
}

func AddResourceByType(w *World, r any, t reflect.Type) bool {
	return hashmap.AddByType(w.resources, r, t)
}

func GetResource[T any](w *World) (*T, bool) {
	return hashmap.Get[T](w.resources)
}

func RemoveResource[T any](w *World) bool {
	return hashmap.Remove[T](w.resources)
}
