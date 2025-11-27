package transport

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

type EventType uint8

const (
	EventConnected EventType = iota
	EventDisconnected
)

type TransportEvent struct {
	Type    EventType
	Address string
}
