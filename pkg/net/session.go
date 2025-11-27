package net

import (
	"github.com/vistormu/xpeto/core/ecs"
)

type session struct {
	lookup   map[string]ecs.Entity
	channels []Channel
}

func updateSession(w *ecs.World) {
	sessions, _ := ecs.GetResource[session](w)

	clear(sessions.lookup)

	q := ecs.NewQuery1[Connection](w)
	for _, b := range q.Iter() {
		conn := b.Components()
		sessions.lookup[conn.Address] = b.Entity()
	}
}
