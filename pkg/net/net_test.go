package net

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
	"net"
	"testing"
	"time"

	"github.com/vistormu/xpeto/core"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/event"
	"github.com/vistormu/xpeto/core/schedule"
)

// ====
// mock
// ====
type Message struct {
	Content string
}

type channels struct {
	UdpGob Channel `protocol:"udp" codec:"gob" listen:"localhost:9876"`
	TcpGob Channel `protocol:"tcp" codec:"gob" listen:"localhost:9877"`
}

const (
	udpAddr = "localhost:9876"
	tcpAddr = "localhost:9877"
)

// =======
// helpers
// =======
func cleanup(t *testing.T, w *ecs.World) func() {
	return func() {
		t.Log("cleaning up network listeners")

		session, ok := ecs.GetResource[session](w)
		if !ok {
			return
		}

		for _, ch := range session.channels {
			if ch.Transport != nil {
				ch.Transport.Close()
			}
		}
	}
}

func setupWorld(t *testing.T) (*ecs.World, *schedule.Scheduler) {
	w := ecs.NewWorld()
	sch := schedule.NewScheduler()

	core.CorePkgs(w, sch)
	Pkg(w, sch)
	AddChannel[channels](w)
	AddMessage[Message](sch, "test_msg")

	schedule.RunStartup(w, sch)
	schedule.RunUpdate(w, sch)
	t.Cleanup(cleanup(t, w))

	return w, sch
}

// =========
// udp + gob
// =========
func TestUdpGob(t *testing.T) {
	w, sch := setupWorld(t)

	// create client
	t.Log("creating upt client")
	raddr, _ := net.ResolveUDPAddr("udp", udpAddr)
	client, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		t.Fatalf("Failed to create UDP client: %v", err)
	}
	defer client.Close()

	// send message
	t.Log("[Action] Sending 'Hello UDP'...")
	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(Message{Content: "Hello UDP"})
	client.Write(buf.Bytes())

	// 4. Wait for Event
	clientAddress := waitForEvent(t, w, sch, ClientConnected)

	// 5. Create Entity
	t.Logf("[Action] Creating Entity for %s...", clientAddress)
	channelsRes, _ := ecs.GetResource[channels](w)
	e := ecs.AddEntity(w)
	ecs.AddComponent(w, e, Connection{
		Address: clientAddress,
		Channel: channelsRes.UdpGob,
	})
	ecs.AddComponent(w, e, Inbox[Message]{Data: make([]Message, 0)})
	ecs.AddComponent(w, e, Outbox[Message]{Data: make([]Message, 0)})

	// 6. Test Receive (Payload)
	t.Log("[Action] Sending 'Payload UDP'...")
	buf.Reset()
	// IMPORTANT: New Encoder for UDP to reset Gob state
	gob.NewEncoder(&buf).Encode(Message{Content: "Payload UDP"})
	client.Write(buf.Bytes())

	// Wait and Verify Inbox
	verifyInbox(t, w, sch, e, "Payload UDP")

	// 7. Test Send (Reply)
	t.Log("[Action] Queuing Reply...")
	outbox, _ := ecs.GetComponent[Outbox[Message]](w, e)
	outbox.Data = append(outbox.Data, Message{Content: "Reply UDP"})

	schedule.RunUpdate(w, sch) // Flush outbox

	// Verify Client Receive
	client.SetReadDeadline(time.Now().Add(time.Second))
	readBuf := make([]byte, 1024)
	n, _, err := client.ReadFromUDP(readBuf)
	if err != nil {
		t.Fatalf("Client read failed: %v", err)
	}

	var reply Message
	gob.NewDecoder(bytes.NewReader(readBuf[:n])).Decode(&reply)
	if reply.Content != "Reply UDP" {
		t.Errorf("Expected 'Reply UDP', got '%s'", reply.Content)
	}
	t.Log("UDP Round Trip Success!")
}

// =============================================================================
// TCP TEST (New)
// =============================================================================

func TestTcpGob(t *testing.T) {
	// 1. Setup
	w, sch := setupWorld(t)

	// 2. Client Setup (TCP)
	t.Log("[Test Setup] Creating TCP Client...")
	// Retry loop for TCP since listener startup might race
	var client net.Conn
	var err error
	for i := 0; i < 5; i++ {
		client, err = net.Dial("tcp", tcpAddr)
		if err == nil {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if err != nil {
		t.Fatalf("Failed to connect to TCP server: %v", err)
	}
	defer client.Close()

	// 3. Connect
	// TCP fires 'Connected' event immediately upon Dial, we don't need to send data first.
	clientAddress := waitForEvent(t, w, sch, ClientConnected)

	// 4. Create Entity
	t.Logf("[Action] Creating Entity for %s...", clientAddress)
	channelsRes, _ := ecs.GetResource[channels](w)
	e := ecs.AddEntity(w)
	ecs.AddComponent(w, e, Connection{
		Address: clientAddress,
		Channel: channelsRes.TcpGob,
	})
	ecs.AddComponent(w, e, Inbox[Message]{Data: make([]Message, 0)})
	ecs.AddComponent(w, e, Outbox[Message]{Data: make([]Message, 0)})

	// 5. Test Receive (With Framing)
	t.Log("[Action] Sending Framed 'Payload TCP'...")
	sendFramedTcp(t, client, Message{Content: "Payload TCP"})

	// Wait and Verify Inbox
	verifyInbox(t, w, sch, e, "Payload TCP")

	// 6. Test Send (Reply)
	t.Log("[Action] Queuing Reply...")
	outbox, _ := ecs.GetComponent[Outbox[Message]](w, e)
	outbox.Data = append(outbox.Data, Message{Content: "Reply TCP"})

	schedule.RunUpdate(w, sch) // Flush outbox

	// Verify Client Receive (Read Frame)
	reply := readFramedTcp(t, client)
	if reply.Content != "Reply TCP" {
		t.Errorf("Expected 'Reply TCP', got '%s'", reply.Content)
	}
	t.Log("TCP Round Trip Success!")
}

func waitForEvent(t *testing.T, w *ecs.World, sch *schedule.Scheduler, targetType ClientEventType) string {
	t.Logf("Waiting for Event %d...", targetType)
	for i := 0; i < 20; i++ { // 20 frames timeout
		time.Sleep(10 * time.Millisecond)
		schedule.RunUpdate(w, sch) // Flush outbox

		events, _ := event.GetEvents[ClientEvent](w)
		for _, e := range events {
			if e.Type == targetType {
				t.Logf("Found Event from %s", e.Address)
				return e.Address
			}
		}
	}
	t.Fatal("Timeout waiting for Client Event")
	return ""
}

func verifyInbox(t *testing.T, w *ecs.World, sch *schedule.Scheduler, e ecs.Entity, expected string) {
	for i := 0; i < 20; i++ {
		time.Sleep(10 * time.Millisecond)
		schedule.RunUpdate(w, sch) // Flush outbox

		inbox, ok := ecs.GetComponent[Inbox[Message]](w, e)
		if ok && len(inbox.Data) > 0 {
			if inbox.Data[0].Content == expected {
				inbox.Data = inbox.Data[:0] // Clear
				return
			}
		}
	}
	t.Fatalf("Timeout waiting for Inbox message: %s", expected)
}

// === TCP FRAMING HELPERS ===

func sendFramedTcp(t *testing.T, conn net.Conn, msg Message) {
	// 1. Encode Gob
	var buf bytes.Buffer
	gob.NewEncoder(&buf).Encode(msg)
	payload := buf.Bytes()

	// 2. Write Header (Uint32 Length)
	header := make([]byte, 4)
	binary.LittleEndian.PutUint32(header, uint32(len(payload)))
	if _, err := conn.Write(header); err != nil {
		t.Fatal(err)
	}

	// 3. Write Body
	if _, err := conn.Write(payload); err != nil {
		t.Fatal(err)
	}
}

func readFramedTcp(t *testing.T, conn net.Conn) Message {
	// 1. Read Header
	header := make([]byte, 4)
	if _, err := io.ReadFull(conn, header); err != nil {
		t.Fatal(err)
	}
	size := binary.LittleEndian.Uint32(header)

	// 2. Read Body
	payload := make([]byte, size)
	if _, err := io.ReadFull(conn, payload); err != nil {
		t.Fatal(err)
	}

	// 3. Decode
	var msg Message
	if err := gob.NewDecoder(bytes.NewReader(payload)).Decode(&msg); err != nil {
		t.Fatal(err)
	}
	return msg
}
