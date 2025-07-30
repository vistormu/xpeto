package event

import (
	"testing"
)

type mockEvent struct {
	data string
}

func TestEventBus(t *testing.T) {
	eb := NewBus()

	// subscribe to a mock event
	value := ""
	Subscribe(eb, func(event mockEvent) {
		value = event.data
	})

	// publish an event
	Publish(eb, mockEvent{data: "Hello, World!"})

	if value != "Hello, World!" {
		t.Errorf("Expected value to be 'Hello, World!', got '%s'", value)
	}
}
