package event

import (
	"sync"
)

// ======
// events
// ======
type entry[T any] struct {
	id   uint64
	data T
}

type events[T any] struct {
	mu       sync.RWMutex
	current  []entry[T]
	previous []entry[T]
	nextId   uint64
}

func newEvents[T any]() events[T] {
	return events[T]{
		current:  make([]entry[T], 0),
		previous: make([]entry[T], 0),
		nextId:   1,
	}
}

func (ev *events[T]) add(e T) uint64 {
	ev.mu.Lock()
	defer ev.mu.Unlock()

	id := ev.nextId
	ev.nextId++

	ev.current = append(ev.current, entry[T]{id: id, data: e})

	return id
}

func (ev *events[T]) update() {
	ev.mu.Lock()
	defer ev.mu.Unlock()

	ev.previous, ev.current = ev.current, ev.previous[:0]
}

// ======
// reader
// ======
type eventReader[T any] struct {
	index uint64
}

func (er *eventReader[T]) read(ev *events[T]) ([]T, bool) {
	ev.mu.RLock()

	count := 0
	for _, e := range ev.previous {
		if e.id > er.index {
			count++
		}
	}
	for _, e := range ev.current {
		if e.id > er.index {
			count++
		}
	}

	if count == 0 {
		ev.mu.RUnlock()
		return nil, false
	}

	out := make([]T, 0, count)
	latest := er.index

	for i := range ev.previous {
		if ev.previous[i].id > er.index {
			out = append(out, ev.previous[i].data)
			if ev.previous[i].id > latest {
				latest = ev.previous[i].id
			}
		}
	}
	for i := range ev.current {
		if ev.current[i].id > er.index {
			out = append(out, ev.current[i].data)
			if ev.current[i].id > latest {
				latest = ev.current[i].id
			}
		}
	}

	ev.mu.RUnlock()

	er.index = latest

	return out, true
}
