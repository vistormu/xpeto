package transport

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type UDPTransport struct {
	// sync
	wg sync.WaitGroup
	mu sync.Mutex

	// conn
	conn   *net.UDPConn
	closed chan struct{}

	// messages
	inbox   chan Packet
	dropped int

	// errors
	errbox     chan error
	errdropped int

	// peers
	peers    map[string]time.Time
	eventBox chan TransportEvent
	timeout  time.Duration

	bufSize int
}

func NewUDP() *UDPTransport {
	return &UDPTransport{
		closed:   make(chan struct{}),
		inbox:    make(chan Packet, 1024),
		errbox:   make(chan error, 100),
		peers:    make(map[string]time.Time),
		eventBox: make(chan TransportEvent, 100),
		timeout:  time.Second * 5,
		bufSize:  4096,
	}
}

func (u *UDPTransport) Listen(address string) error {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	u.conn = conn

	u.wg.Add(1)
	go u.readLoop()

	return nil
}

func (u *UDPTransport) readLoop() {
	defer u.wg.Done()

	buf := make([]byte, u.bufSize)

	for {
		select {
		case <-u.closed:
			return
		default:
		}

		n, addr, err := u.conn.ReadFromUDP(buf)
		if err != nil {
			err = fmt.Errorf("udp read error. %w", err)

			select {
			case u.errbox <- err:
			default:
				u.errdropped++
			}
			continue
		}

		sender := addr.String()
		u.mu.Lock()
		_, ok := u.peers[sender]
		u.peers[sender] = time.Now()
		u.mu.Unlock()

		if !ok {
			select {
			case u.eventBox <- TransportEvent{Type: EventConnected, Address: sender}:
			default:
			}
		}

		data := make([]byte, n)
		copy(data, buf[:n])

		pkt := Packet{
			Sender:  addr.String(),
			Payload: data,
		}

		select {
		case u.inbox <- pkt:
		default:
			u.dropped++
		}
	}
}

func (u *UDPTransport) Send(to string, data []byte) error {
	if u.conn == nil {
		return net.ErrClosed
	}

	addr, err := net.ResolveUDPAddr("udp", to)
	if err != nil {
		return err
	}

	_, err = u.conn.WriteToUDP(data, addr)

	return err
}

func (u *UDPTransport) Flush(max int) []Packet {
	n := len(u.inbox)
	if n == 0 {
		return nil
	}

	n = min(max, n)

	packets := make([]Packet, 0, n)

loop:
	for range n {
		select {
		case pkt := <-u.inbox:
			packets = append(packets, pkt)
		default:
			break loop
		}
	}

	return packets
}

func (u *UDPTransport) FlushErrors() []error {
	n := len(u.errbox)
	if n == 0 {
		return nil
	}

	errors := make([]error, 0, n)
loop:
	for {
		select {
		case err := <-u.errbox:
			errors = append(errors, err)
		default:
			break loop
		}
	}

	return errors
}

func (u *UDPTransport) FlushEvents() []TransportEvent {
	var events []TransportEvent

loop:
	for {
		select {
		case ev := <-u.eventBox:
			events = append(events, ev)
		default:
			break loop
		}
	}

	now := time.Now()
	u.mu.Lock()
	defer u.mu.Unlock()

	for addr, lastSeen := range u.peers {
		if now.Sub(lastSeen) > u.timeout {
			events = append(events, TransportEvent{
				Type:    EventDisconnected,
				Address: addr,
			})
			delete(u.peers, addr)
		}
	}

	return events
}

func (u *UDPTransport) Close() error {
	close(u.closed)

	var err error
	if u.conn != nil {
		err = u.conn.Close()
	}

	u.wg.Wait()

	close(u.inbox)

	return err
}
