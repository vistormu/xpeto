package event

import (
	"reflect"
	"slices"
)

type Manager struct {
	callbacks map[reflect.Type][]func(event any)
}

func NewManager() *Manager {
	return &Manager{
		callbacks: make(map[reflect.Type][]func(event any)),
	}
}

func Subscribe[T any](manager *Manager, callback func(event T)) {
	// get the type of the event
	eventType := reflect.TypeOf((*T)(nil)).Elem()

	// check if the event type is already registered
	if _, ok := manager.callbacks[eventType]; !ok {
		manager.callbacks[eventType] = []func(event any){}
	}

	// add the callback to the list of callbacks for the event type
	manager.callbacks[eventType] = append(manager.callbacks[eventType], func(event any) {
		callback(event.(T))
	})
}

func Unsubscribe[T any](manager *Manager, callback func(event T)) {
	// get the type of the event
	eventType := reflect.TypeOf((*T)(nil)).Elem()

	// check if the event type is registered
	callbacks, ok := manager.callbacks[eventType]
	if !ok {
		return
	}

	// remove the callback from the list of callbacks for the event type
	for i, cb := range callbacks {
		if reflect.ValueOf(cb).Pointer() == reflect.ValueOf(callback).Pointer() {
			manager.callbacks[eventType] = slices.Delete(callbacks, i, i+1)
			break
		}
	}
}

func (m *Manager) Publish(event any) {
	// get the type of the event
	eventType := reflect.TypeOf(event)

	// check if there are any subscribers to the event
	callbacks, ok := m.callbacks[eventType]
	if !ok {
		return
	}

	// call the callbacks
	for _, callback := range callbacks {
		callback(event)
	}
}
