package ecs

import (
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

func newStore[T any]() *store[T] {
	return &store[T]{
		dense:    make([]Entity, 0),
		values:   make([]T, 0),
		location: make(map[Entity]int),
	}
}

func (s *store[T]) add(e Entity, v T) {
	row, ok := s.location[e]
	if ok {
		s.values[row] = v
		return
	}

	row = len(s.dense)
	s.dense = append(s.dense, e)
	s.values = append(s.values, v)
	s.location[e] = row
}

func (s *store[T]) get(e Entity) (*T, bool) {
	row, ok := s.location[e]
	if !ok {
		return nil, false
	}
	return &s.values[row], true
}

func (s *store[T]) remove(e Entity) bool {
	row, ok := s.location[e]
	if !ok {
		return false
	}

	last := len(s.dense) - 1
	if row != last {
		le := s.dense[last]
		s.dense[row] = le
		s.values[row] = s.values[last]
		s.location[le] = row
	}

	s.dense = s.dense[:last]
	s.values = s.values[:last]

	delete(s.location, e)

	return true
}

// ========
// registry
// ========
func getStore[T any](r *hashmap.TypeMap) *store[T] {
	s, ok := hashmap.Get[store[T]](r)
	if ok {
		return s
	}

	ns := newStore[T]()
	hashmap.Add(r, ns)

	return ns
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

// =================
// helper components
// =================
type Name struct {
	Value string
}
