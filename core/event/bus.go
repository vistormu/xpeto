package event

import (
	"reflect"
	"slices"
)

type Bus struct {
	events    map[reflect.Type][]Event
	callbacks map[Event]func(data any)
	nextId    uint32
}

func NewBus() *Bus {
	return &Bus{
		events:    make(map[reflect.Type][]Event),
		callbacks: make(map[Event]func(data any)),
		nextId:    1,
	}
}

func Subscribe[T any](b *Bus, callback func(data T)) Event {
	// register callback
	eventType := reflect.TypeFor[T]()
	event := Event{Id: b.nextId}
	b.nextId++

	// create new callback slice if it doesn't exist
	_, ok := b.events[eventType]
	if !ok {
		b.events[eventType] = []Event{}
	}

	// add the callbacks to the map and events
	b.callbacks[event] = func(event any) {
		callback(event.(T))
	}

	b.events[eventType] = append(b.events[eventType], event)

	return event
}

func (b *Bus) Unsubscribe(event Event) {
	// remove the callback from the map
	_, ok := b.callbacks[event]
	if !ok {
		return
	}
	delete(b.callbacks, event)

	// remove the event from all event types
	for eventType, events := range b.events {
		for i, h := range events {
			if h.Id == event.Id {
				b.events[eventType] = slices.Delete(events, i, i+1)
				break
			}
		}
	}
}

func Publish[T any](b *Bus, data T) {
	// get the type of the event
	eventType := reflect.TypeFor[T]()

	// check if there are any subscribers to the event
	events, ok := b.events[eventType]
	if !ok {
		return
	}

	// call the events
	for _, e := range events {
		// get the callback function
		callback, ok := b.callbacks[e]
		if !ok {
			continue
		}

		callback(data)
	}
}
