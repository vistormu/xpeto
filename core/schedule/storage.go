package schedule

import (
	"github.com/vistormu/go-dsa/hashmap"
	"github.com/vistormu/go-dsa/set"
)

type storage struct {
	nodes         []*node
	stages        map[stage][]uint64
	plan          map[stage][]uint64
	dirty         *set.HashSet[stage]
	stateMachines *hashmap.TypeMap
}

func newStorage() *storage {
	return &storage{
		nodes:         make([]*node, 0),
		stages:        make(map[stage][]uint64),
		plan:          make(map[stage][]uint64),
		dirty:         set.NewHashSet[stage](),
		stateMachines: hashmap.NewTypeMap(),
	}
}

func (s *storage) add(n *node) {
	s.nodes = append(s.nodes, n)

	if n.stage == empty {
		return
	}

	s.stages[n.stage] = append(s.stages[n.stage], n.id)
	s.dirty.Add(n.stage)
}

func (s *storage) get(id uint64) (*node, bool) {
	if id == 0 {
		return nil, false
	}

	i := int(id - 1)
	if i < 0 || i >= len(s.nodes) {
		return nil, false
	}

	return s.nodes[i], true
}

func addStateMachine[T comparable](store *storage, sm stateMachine[T]) {
	hashmap.Add(store.stateMachines, sm)
}

func getStateMachine[T comparable](store *storage) (*stateMachine[T], bool) {
	sm, ok := hashmap.Get[stateMachine[T]](store.stateMachines)
	return sm, ok
}
