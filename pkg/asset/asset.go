package asset

import "github.com/vistormu/go-dsa/queue"

// =====
// asset
// =====
type Asset uint64

func newAsset(index, gen uint32) Asset {
	return Asset(uint64(gen)<<32 | uint64(index))
}

func (a Asset) gen() uint32 {
	return uint32(uint64(a) >> 32)
}

func (a Asset) index() uint32 {
	return uint32(uint64(a) & 0xffffffff)
}

// ==========
// population
// ==========
type population struct {
	free  *queue.QueueArray[uint32]
	gens  []uint32
	alive uint32
}

func newPopulation() *population {
	return &population{
		free:  queue.NewQueueArray[uint32](),
		gens:  make([]uint32, 0),
		alive: 0,
	}
}

func (p *population) add() Asset {
	var index uint32

	if !p.free.IsEmpty() {
		index, _ = p.free.Dequeue()
	} else {
		index = uint32(len(p.gens))
		p.gens = append(p.gens, 1)
	}

	gen := p.gens[index]
	p.alive++

	return newAsset(index, gen)
}

func (p *population) remove(a Asset) bool {
	index := a.index()

	if index >= uint32(len(p.gens)) || p.gens[index] != a.gen() {
		return false
	}

	p.gens[index]++
	p.free.Enqueue(index)

	p.alive--

	return true
}

// func (p *population) has(a Asset) bool {
// 	index := a.index()

// 	return index < uint32(len(p.gens)) && p.gens[index] == a.gen()
// }
