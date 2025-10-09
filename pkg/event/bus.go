package event

import (
	"github.com/vistormu/go-dsa/hashmap"
	"github.com/vistormu/xpeto/core/ecs"
)

type bus struct {
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
	eb, _ := ecs.GetResource[bus](w)

	for _, e := range eb.registry.Values() {
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
	eb, _ := ecs.GetResource[bus](w)

	_, ok := hashmap.Get[events[T]](eb.registry)
	if !ok {
		hashmap.Add(eb.registry, newEvents[T]())
	}

	ev, _ := hashmap.Get[events[T]](eb.registry)
	ev.add(e)
}

func GetEvents[T any](w *ecs.World) ([]T, bool) {
	eb, _ := ecs.GetResource[bus](w)
	id := ecs.GetSystemId(w)

	rm, ok := eb.readers[id]
	if !ok {
		rm = hashmap.NewTypeMap()
		eb.readers[id] = rm
	}

	er, ok := hashmap.Get[eventReader[T]](rm)
	if !ok {
		hashmap.Add(rm, eventReader[T]{})
		er, _ = hashmap.Get[eventReader[T]](rm)
	}

	ev, ok := hashmap.Get[events[T]](eb.registry)
	if !ok {
		return nil, false
	}

	return er.read(ev)
}
