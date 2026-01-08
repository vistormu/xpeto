package ecs

import (
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

	return w
}

// ===
// API
// ===
func AddResource[T any](w *World, r T) bool {
	return hashmap.Add(w.resources, r)
}

func EnsureResource[T any](w *World, init func() T) *T {
	r, ok := GetResource[T](w)
	if ok {
		return r
	}
	AddResource(w, init())
	r, _ = GetResource[T](w)
	return r
}

func RemoveResource[T any](w *World) bool {
	return hashmap.Remove[T](w.resources)
}

func GetResource[T any](w *World) (*T, bool) {
	return hashmap.Get[T](w.resources)
}

func HasResource[T any](w *World) bool {
	return hashmap.Has[T](w.resources)
}
