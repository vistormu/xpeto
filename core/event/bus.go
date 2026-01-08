package event

import (
	"sync"

	"github.com/vistormu/go-dsa/hashmap"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

type bus struct {
	mu       sync.RWMutex
	registry *hashmap.TypeMap
	readers  map[uint64]*hashmap.TypeMap
}

func newBus() bus {
	return bus{
		registry: hashmap.NewTypeMap(),
		readers:  make(map[uint64]*hashmap.TypeMap),
	}
}

func update(w *ecs.World) {
	eb := ecs.EnsureResource(w, newBus)

	eb.mu.RLock()
	values := eb.registry.Values()
	eb.mu.RUnlock()

	for _, e := range values {
		switch ev := e.(type) {
		case interface{ update() }:
			ev.update()
		}
	}
}

// ===
// API
// ===
func AddEvent[T any](w *ecs.World, e T) {
	eb := ecs.EnsureResource(w, newBus)
	eb.mu.Lock()
	defer eb.mu.Unlock()

	_, ok := hashmap.Get[events[T]](eb.registry)
	if !ok {
		hashmap.Add(eb.registry, newEvents[T]())
	}

	ev, _ := hashmap.Get[events[T]](eb.registry)
	ev.add(e)
}

func GetEvents[T any](w *ecs.World) ([]T, bool) {
	eb, ok := ecs.GetResource[bus](w)
	if !ok {
		return nil, false
	}
	rs, ok := ecs.GetResource[schedule.RunningSystem](w)
	if !ok {
		return nil, false
	}

	eb.mu.Lock()
	rm, ok := eb.readers[rs.Id]
	if !ok {
		rm = hashmap.NewTypeMap()
		eb.readers[rs.Id] = rm
	}

	er, ok := hashmap.Get[eventReader[T]](rm)
	if !ok {
		hashmap.Add(rm, eventReader[T]{})
		er, _ = hashmap.Get[eventReader[T]](rm)
	}

	ev, ok := hashmap.Get[events[T]](eb.registry)
	eb.mu.Unlock()

	if !ok {
		return nil, false
	}

	return er.read(ev)
}
