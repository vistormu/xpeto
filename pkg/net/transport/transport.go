package transport

import (
	"strings"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/log"
)

type Packet struct {
	Sender  string
	Payload []byte
}

type Transport interface {
	Listen(address string) error
	Send(to string, data []byte) error
	Flush(max int) []Packet
	FlushErrors() []error
	FlushEvents() []TransportEvent
	Close() error
}

func New(w *ecs.World, c string) Transport {
	switch strings.ToLower(c) {
	case "udp":
		return newUDP()

	case "tcp":
		return newTCP()

	default:
		log.LogError(w, "unknown protocol", log.F("got", c))
		return nil
	}
}

type EventType uint8

const (
	EventConnected EventType = iota
	EventDisconnected
)

type TransportEvent struct {
	Type    EventType
	Address string
}
