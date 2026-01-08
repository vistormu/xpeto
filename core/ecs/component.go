package ecs

import (
	"reflect"

	"github.com/vistormu/go-dsa/hashmap"
)

// =====
// store
// =====
type store[T any] struct {
	dense    []Entity
	values   []T
	location map[Entity]int
}

func newStore[T any]() store[T] {
	return store[T]{
		dense:    make([]Entity, 0),
		values:   make([]T, 0),
		location: make(map[Entity]int),
	}
}

func (s *store[T]) add(e Entity, c T) {
	row, ok := s.location[e]
	if ok {
		s.values[row] = c
		return
	}

	s.location[e] = len(s.dense)
	s.dense = append(s.dense, e)
	s.values = append(s.values, c)
}

func (s *store[T]) remove(e Entity) bool {
	row, ok := s.location[e]
	if !ok {
		return false
	}

	last := len(s.dense) - 1
	lastE := s.dense[last]

	s.dense[row] = s.dense[last]
	s.values[row] = s.values[last]

	s.dense = s.dense[:last]
	s.values = s.values[:last]

	delete(s.location, e)
	if row != last {
		s.location[lastE] = row
	}

	return true
}

func (s *store[T]) get(e Entity) (*T, bool) {
	row, ok := s.location[e]
	if !ok {
		return nil, false
	}
	return &s.values[row], true
}

func (s *store[T]) has(e Entity) bool {
	_, ok := s.location[e]
	return ok
}

// ========
// registry
// ========
func getStore[T any](r *hashmap.TypeMap) *store[T] {
	_, ok := hashmap.Get[store[T]](r)
	if !ok {
		hashmap.Add(r, newStore[T]())
	}

	s, _ := hashmap.Get[store[T]](r)

	return s
}

// TODO: the complexity is O(c), where c is the number of components
func removeComponents(r *hashmap.TypeMap, e Entity) bool {
	removed := false
	for _, raw := range r.Iter() {
		switch s := raw.(type) {
		case interface{ remove(Entity) bool }:
			if s.remove(e) {
				removed = true
			}
		}
	}

	return removed
}

// ===
// API
// ===
func AddComponent[T any](w *World, e Entity, c T) bool {
	if !w.population.has(e) {
		return false
	}

	if reflect.TypeFor[T]().Kind() == reflect.Pointer {
		return false
	}

	getStore[T](w.registry).add(e, c)

	return true
}

func RemoveComponent[T any](w *World, e Entity) bool {
	if !w.population.has(e) {
		return false
	}

	return getStore[T](w.registry).remove(e)
}

func GetComponent[T any](w *World, e Entity) (*T, bool) {
	if !w.population.has(e) {
		return nil, false
	}

	return getStore[T](w.registry).get(e)
}

func HasComponent[T any](w *World, e Entity) bool {
	if !w.population.has(e) {
		return false
	}
	return getStore[T](w.registry).has(e)
}

// func ReserveComponents[T any](w *World, n int) {
// 	if n <= 0 {
// 		return
// 	}

// 	s := getStore[T](w.registry)

// 	if cap(s.dense) < n {
// 		d := make([]Entity, len(s.dense), n)
// 		copy(d, s.dense)
// 		s.dense = d
// 	}

// 	if cap(s.values) < n {
// 		v := make([]T, len(s.values), n)
// 		copy(v, s.values)
// 		s.values = v
// 	}

// 	if s.location == nil {
// 		s.location = make(map[Entity]int, n)
// 	}
// }
