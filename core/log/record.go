package log

import (
	"fmt"
	"time"
)

type Record struct {
	Level       Level
	SystemId    uint64
	SystemLabel string
	Frame       uint64
	Time        time.Duration
	Caller      Caller
	Message     string
	Fields      []field
}

type field struct {
	key   string
	value any
}

func (f field) String() string {
	return fmt.Sprintf("%s: %v", f.key, f.value)
}

func (f field) Key() string {
	return f.key
}

func (f field) Value() any {
	return f.value
}

func F(key string, value any) field {
	return field{key, value}
}
