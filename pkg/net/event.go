package net

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/event"
	"github.com/vistormu/xpeto/pkg/net/transport"
)

type ClientEventType uint8

const (
	ClientConnected ClientEventType = iota
	ClientDisconnected
)

type ClientEvent struct {
	Type    ClientEventType
	Address string
}

func emitEvents(w *ecs.World) {
	session, _ := ecs.GetResource[session](w)
	if len(session.channels) == 0 {
		return
	}

	for _, ch := range session.channels {
		for _, e := range ch.Transport.FlushEvents() {
			switch e.Type {
			case transport.EventConnected:
				event.AddEvent(w, ClientEvent{
					Type:    ClientConnected,
					Address: e.Address,
				})
			case transport.EventDisconnected:
				event.AddEvent(w, ClientEvent{
					Type:    ClientDisconnected,
					Address: e.Address,
				})

				en, ok := session.lookup[e.Address]
				if ok {
					delete(session.lookup, e.Address)
					// delete entity with connection?
					ecs.RemoveComponent[Connection](w, en)
				}
			}
		}
	}
}
