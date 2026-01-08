package asset

import (
	"github.com/vistormu/go-dsa/hashmap"
)

// =====
// store
// =====
type store[T any] struct {
	dense    []Asset
	values   []*T
	location map[Asset]int
}

func newStore[T any]() store[T] {
	return store[T]{
		dense:    make([]Asset, 0),
		values:   make([]*T, 0),
		location: make(map[Asset]int),
	}
}

func (s *store[T]) add(a Asset, v *T) {
	row, ok := s.location[a]
	if ok {
		s.values[row] = v
		return
	}

	row = len(s.dense)
	s.dense = append(s.dense, a)
	s.values = append(s.values, v)
	s.location[a] = row
}

func (s *store[T]) get(a Asset) (*T, bool) {
	row, ok := s.location[a]
	if !ok {
		return nil, false
	}
	return s.values[row], true
}

func (s *store[T]) remove(a Asset) bool {
	row, ok := s.location[a]
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

	delete(s.location, a)

	return true
}

func getStore[T any](r *hashmap.TypeMap) *store[T] {
	_, ok := hashmap.Get[store[T]](r)
	if !ok {
		hashmap.Add(r, newStore[T]())
	}

	s, _ := hashmap.Get[store[T]](r)

	return s
}
