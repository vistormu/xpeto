package ecs

import (
	"github.com/vistormu/go-dsa/queue"
)

// ======
// entity
// ======
type Entity uint64

func newEntity(index, gen uint32) Entity {
	return Entity(uint64(gen)<<32 | uint64(index))
}

func (e Entity) gen() uint32 {
	return uint32(uint64(e) >> 32)
}

func (e Entity) index() uint32 {
	return uint32(uint64(e) & 0xffffffff)
}

// ==========
// population
// ==========
type population struct {
	free  *queue.Queue[uint32]
	gens  []uint32
	alive uint32
}

func newPopulation() *population {
	return &population{
		free:  queue.NewQueue[uint32](),
		gens:  make([]uint32, 0),
		alive: 0,
	}
}

func (p *population) add() Entity {
	var index uint32

	if !p.free.Empty() {
		index, _ = p.free.Dequeue()
	} else {
		index = uint32(len(p.gens))
		p.gens = append(p.gens, 1)
	}

	gen := p.gens[index]
	p.alive++

	return newEntity(index, gen)
}

func (p *population) remove(e Entity) bool {
	index := e.index()

	if index >= uint32(len(p.gens)) || p.gens[index] != e.gen() {
		return false
	}

	p.gens[index]++
	p.free.Enqueue(index)

	p.alive--

	return true
}

func (p *population) has(e Entity) bool {
	index := e.index()

	return index < uint32(len(p.gens)) && p.gens[index] == e.gen()
}

// ===
// API
// ===
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
