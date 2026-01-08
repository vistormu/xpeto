package asset

import (
	"github.com/vistormu/go-dsa/queue"
)

// =====
// state
// =====
type AssetState uint8

const (
	AssetNone AssetState = iota
	AssetRequested
	AssetLoading
	AssetLoaded
	AssetFailed
)

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
	free  *queue.Queue[uint32]
	gens  []uint32
	alive uint32

	state []AssetState
	path  []string
	err   []error
}

func newPopulation() *population {
	return &population{
		free:  queue.NewQueue[uint32](),
		gens:  make([]uint32, 0),
		alive: 0,
		state: make([]AssetState, 0),
		path:  make([]string, 0),
		err:   make([]error, 0),
	}
}

func (p *population) add() Asset {
	var index uint32

	if !p.free.Empty() {
		index, _ = p.free.Dequeue()
	} else {
		index = uint32(len(p.gens))
		p.gens = append(p.gens, 1)
		p.state = append(p.state, AssetNone)
		p.path = append(p.path, "")
		p.err = append(p.err, nil)
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

	p.state[index] = AssetNone
	p.path[index] = ""
	p.err[index] = nil

	return true
}

func (p *population) has(a Asset) bool {
	i := a.index()
	return i < uint32(len(p.gens)) && p.gens[i] == a.gen()
}

func (p *population) setRequested(a Asset, path string) bool {
	if !p.has(a) {
		return false
	}
	i := a.index()
	p.state[i] = AssetRequested
	p.path[i] = path
	p.err[i] = nil
	return true
}

func (p *population) setLoading(a Asset) bool {
	if !p.has(a) {
		return false
	}
	p.state[a.index()] = AssetLoading
	return true
}

func (p *population) setLoaded(a Asset) bool {
	if !p.has(a) {
		return false
	}
	i := a.index()
	p.state[i] = AssetLoaded
	p.err[i] = nil
	return true
}

func (p *population) setFailed(a Asset, err error) bool {
	if !p.has(a) {
		return false
	}
	i := a.index()
	p.state[i] = AssetFailed
	p.err[i] = err
	return true
}

func (p *population) getState(a Asset) (AssetState, bool) {
	if !p.has(a) {
		return AssetNone, false
	}
	return p.state[a.index()], true
}

// func (p *population) isLoaded(a Asset) bool {
// 	st, ok := p.getState(a)
// 	return ok && st == AssetLoaded
// }
