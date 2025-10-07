package ecs

import (
	"reflect"
)

type World struct {
	population *population
	registry   *registry
	resources  map[reflect.Type]any
}

func NewWorld() *World {
	return &World{
		population: newPopulation(),
		registry:   newRegistry(),
		resources:  make(map[reflect.Type]any),
	}
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

	w.registry.removeComponents(e)

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
func baseType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Pointer {
		return t.Elem()
	}
	return t
}

func asPointerToValue(v any) any {
	rv := reflect.ValueOf(v)
	rt := rv.Type()

	if rt.Kind() == reflect.Pointer {
		return v
	}

	ptr := reflect.New(rt)
	ptr.Elem().Set(rv)

	return ptr.Interface()
}

func AddResource[T any](w *World, r T) {
	rt := reflect.TypeFor[T]()
	key := baseType(rt)
	ptr := asPointerToValue(r)
	w.resources[key] = ptr
}

func AddResourceByType(w *World, r any, t reflect.Type) bool {
	if t == nil || r == nil {
		return false
	}

	key := baseType(t)

	rv := reflect.ValueOf(r)
	rt := rv.Type()

	switch {
	case rt == t || rt.AssignableTo(t):
		ptr := reflect.New(key)
		ptr.Elem().Set(rv.Convert(key))
		w.resources[key] = ptr.Interface()

		return true

	case rt == reflect.PointerTo(key) || rt.AssignableTo(reflect.PointerTo(key)):
		w.resources[key] = r

		return true

	default:
		return false
	}
}

func GetResource[T any](w *World) (*T, bool) {
	rt := reflect.TypeFor[T]()
	if rt.Kind() == reflect.Pointer {
		return nil, false
	}

	v, ok := w.resources[rt]
	if !ok {
		return nil, false
	}

	out, ok := v.(*T)

	return out, ok
}

func RemoveResource[T any](w *World) bool {
	rt := reflect.TypeFor[T]()
	if rt.Kind() == reflect.Pointer {
		return false
	}

	if _, ok := w.resources[rt]; !ok {
		return false
	}

	delete(w.resources, rt)

	return true
}
