package event

import (
	"reflect"
	"slices"
)

type Manager struct {
	events    map[reflect.Type][]Event
	callbacks map[Event]func(event any)
	nextId    uint32
}

func NewManager() *Manager {
	return &Manager{
		events:    make(map[reflect.Type][]Event),
		callbacks: make(map[Event]func(event any)),
		nextId:    1,
	}
}

func Subscribe[T any](manager *Manager, callback func(event T)) Event {
	// register callback
	eventType := reflect.TypeOf((*T)(nil)).Elem()
	event := Event{Id: manager.nextId}
	manager.nextId++

	// create new callback slice if it doesn't exist
	_, ok := manager.events[eventType]
	if !ok {
		manager.events[eventType] = []Event{}
	}

	// add the callbacks to the map and events
	manager.callbacks[event] = func(event any) {
		callback(event.(T))
	}

	manager.events[eventType] = append(manager.events[eventType], event)

	return event
}

func (m *Manager) Unsubscribe(event Event) {
	// remove the callback from the map
	_, ok := m.callbacks[event]
	if !ok {
		return
	}
	delete(m.callbacks, event)

	// remove the event from all event types
	for eventType, events := range m.events {
		for i, h := range events {
			if h.Id == event.Id {
				m.events[eventType] = slices.Delete(events, i, i+1)
				break
			}
		}
	}
}

func (m *Manager) Publish(event any) {
	// get the type of the event
	eventType := reflect.TypeOf(event)

	// check if there are any subscribers to the event
	events, ok := m.events[eventType]
	if !ok {
		return
	}

	// call the events
	for _, event := range events {
		// get the callback function
		callback, ok := m.callbacks[event]
		if !ok {
			continue
		}

		callback(event)
	}
}
