package transport

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
)

type TCPTransport struct {
	mu    sync.RWMutex
	conns map[string]net.Conn

	listener net.Listener
	closed   chan struct{}

	inbox    chan Packet
	errbox   chan error
	eventBox chan TransportEvent

	dropped int
}

func newTCP() *TCPTransport {
	return &TCPTransport{
		conns:    make(map[string]net.Conn),
		closed:   make(chan struct{}),
		inbox:    make(chan Packet, 1024),
		errbox:   make(chan error, 100),
		eventBox: make(chan TransportEvent, 100),
	}
}

func (t *TCPTransport) Listen(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	t.listener = l

	go t.acceptLoop()

	return nil
}

func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			select {
			case <-t.closed:
				return
			default:
				t.errbox <- fmt.Errorf("tcp accept error: %w", err)
				continue
			}
		}

		remoteAddr := conn.RemoteAddr().String()
		t.mu.Lock()
		t.conns[remoteAddr] = conn
		t.mu.Unlock()

		t.eventBox <- TransportEvent{Type: EventConnected, Address: remoteAddr}

		go t.handleConnection(conn, remoteAddr)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn, addr string) {
	defer func() {
		conn.Close()
		t.mu.Lock()
		delete(t.conns, addr)
		t.mu.Unlock()

		t.eventBox <- TransportEvent{Type: EventDisconnected, Address: addr}
	}()

	header := make([]byte, 4)

	for {
		_, err := io.ReadFull(conn, header)
		if err != nil {
			if err != io.EOF {
				t.errbox <- fmt.Errorf("tcp read header error from %s: %w", addr, err)
			}
			return
		}

		size := binary.LittleEndian.Uint32(header)

		payload := make([]byte, size)
		_, err = io.ReadFull(conn, payload)
		if err != nil {
			t.errbox <- fmt.Errorf("tcp read body error from %s: %w", addr, err)
			return
		}

		select {
		case t.inbox <- Packet{Sender: addr, Payload: payload}:
		default:
			t.dropped++
		}
	}
}

func (t *TCPTransport) Send(to string, data []byte) error {
	t.mu.RLock()
	conn, ok := t.conns[to]
	t.mu.RUnlock()

	if !ok {
		var err error
		conn, err = net.Dial("tcp", to)
		if err != nil {
			return err
		}

		t.mu.Lock()
		t.conns[to] = conn
		t.mu.Unlock()

		go t.handleConnection(conn, to)

		t.eventBox <- TransportEvent{Type: EventConnected, Address: to}
	}

	frame := make([]byte, 4+len(data))
	binary.LittleEndian.PutUint32(frame[:4], uint32(len(data)))
	copy(frame[4:], data)

	_, err := conn.Write(frame)
	return err
}

func (t *TCPTransport) Flush(max int) []Packet {
	n := len(t.inbox)
	if n == 0 {
		return nil
	}
	n = min(n, max)

	packets := make([]Packet, 0, n)
	for i := 0; i < n; i++ {
		select {
		case p := <-t.inbox:
			packets = append(packets, p)
		default:
		}
	}

	return packets
}

func (t *TCPTransport) FlushErrors() []error {
	var errs []error
loop:
	for {
		select {
		case e := <-t.errbox:
			errs = append(errs, e)
		default:
			break loop
		}
	}
	return errs
}

func (t *TCPTransport) FlushEvents() []TransportEvent {
	var evts []TransportEvent
loop:
	for {
		select {
		case e := <-t.eventBox:
			evts = append(evts, e)
		default:
			break loop
		}
	}
	return evts
}

func (t *TCPTransport) Close() error {
	close(t.closed)
	if t.listener != nil {
		t.listener.Close()
	}

	t.mu.Lock()
	for _, c := range t.conns {
		c.Close()
	}
	t.mu.Unlock()

	return nil
}
