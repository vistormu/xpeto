package log

import (
	"time"
)

type record struct {
	level       Level
	systemId    uint64
	systemLabel string
	frame       uint64
	time        time.Duration
	message     string
	fields      []field
}

type field struct {
	key   string
	value any
}

func F(key string, value any) field {
	return field{key, value}
}
