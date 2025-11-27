package net

import (
	"bytes"
	"encoding/gob"
	"net"
	"testing"
	"time"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/pkg"
	"github.com/vistormu/xpeto/core/pkg/event"
	"github.com/vistormu/xpeto/core/schedule"
)

// ====
// mock
// ====
type Message struct {
	Content string
}

// Note: Using 0.0.0.0 to ensure we bind to all interfaces (IPv4)
type channels struct {
	UdpGob Channel `protocol:"udp" codec:"gob" listen:"0.0.0.0:9876"`
}

const serverAddr = "127.0.0.1:9876"

// =======
// helpers
// =======
func newUdpClient(t *testing.T) *net.UDPConn {
	// Debug: Print what we are dialing
	t.Logf("[Test Setup] Creating Client dialing %s...", serverAddr)

	addr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		t.Fatalf("Failed to resolve address: %v", err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}

	return conn
}

func TestUdpGob(t *testing.T) {
	// initialize app
	w := ecs.NewWorld()
	sch := schedule.NewScheduler()

	t.Log("[Init] Initializing Core...")
	pkg.CorePkgs(w, sch)

	// initialize net
	t.Log("[Init] Initializing Net...")
	Pkg(w, sch)
	AddChannel[channels](w)
	AddMessage[Message](sch, "test_msg")

	// =========================================================================
	// PROBE A: SESSION INTEGRITY
	// Check if AddChannel correctly populated the internal session resource
	// =========================================================================
	sess, ok := ecs.GetResource[session](w)
	if !ok {
		t.Fatal("[PROBE A] FAIL: Session resource not found in World")
	}
	if len(sess.channels) == 0 {
		t.Fatal("[PROBE A] FAIL: Session has 0 channels. AddChannel failed to append.")
	}
	t.Logf("[PROBE A] OK: Session has %d active channel(s)", len(sess.channels))

	// =========================================================================
	// PROBE B: OS LISTENER CHECK
	// Check if the OS actually opened the port
	// =========================================================================
	connCheck, err := net.DialTimeout("udp", serverAddr, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("[PROBE B] FAIL: Could not reach own server port %s: %v", serverAddr, err)
	}
	connCheck.Close()
	t.Logf("[PROBE B] OK: Port %s is open and accepting traffic", serverAddr)

	// Setup Client
	client := newUdpClient(t)
	defer client.Close()

	// Run Startup
	sch.RunStartup(w)

	// Encode message
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err = encoder.Encode(Message{Content: "Hello Server"})
	if err != nil {
		t.Fatal(err)
	}

	// Send Packet
	t.Log("[Action] Sending 'Hello Server' packet...")
	_, err = client.Write(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	// Step engine until client connected
	var clientAddress string
	foundConnect := false

	t.Log("[Loop] Starting Game Loop to catch event...")

loop:
	for i := range 10 {
		// Wait for OS buffer
		time.Sleep(time.Millisecond * 20)

		// Run ONE frame
		sch.RunUpdate(w)

		// =====================================================================
		// PROBE C: CACHE INSPECTION
		// Did 'dispatch' pick up the packet from the Transport?
		// =====================================================================
		// We inspect the internal cache resource directly
		c, _ := ecs.GetResource[cache](w)
		t.Logf("   Frame %d: Cache contains %d packets", i, len(c.packets))

		// If cache has packets, print their sender to verify
		for pIdx, p := range c.packets {
			t.Logf("      -> Packet[%d] from: %s (Size: %d)", pIdx, p.Sender, len(p.Payload))
		}

		// =====================================================================
		// PROBE D: EVENT INSPECTION
		// Did 'emitEvents' find the event?
		// =====================================================================
		events, _ := event.GetEvents[ClientEvent](w)
		if len(events) > 0 {
			t.Logf("   Frame %d: Found %d ClientEvents in World", i, len(events))
		}

		for _, e := range events {
			t.Logf("      -> Event Type: %d, Addr: %s", e.Type, e.Address)
			if e.Type == ClientConnected {
				foundConnect = true
				clientAddress = e.Address
				t.Logf("[SUCCESS] Client connected from: %s in %d frames", clientAddress, i)
				break loop
			}
		}
	}

	if !foundConnect {
		// Final Debug Dump before failing
		t.Log("---------------------------------------------------")
		t.Log("DEBUG SUMMARY:")
		t.Log("1. If Cache was always 0: 'dispatch' isn't flushing the transport.")
		t.Log("   (Possible causes: IPv4/IPv6 mismatch, or Transport.readLoop not running)")
		t.Log("2. If Cache had packets but No Event: 'emitEvents' didn't trigger.")
		t.Log("   (Possible causes: Packet sender already in session lookup?)")
		t.Log("---------------------------------------------------")
		t.Fatal("Did not find ClientConnected event")
	}

	// =========================================================================
	// PHASE 2: DATA TRANSFER
	// =========================================================================

	// Add receiver and sender
	channelsRes, _ := ecs.GetResource[channels](w)
	e := ecs.AddEntity(w)
	ecs.AddComponent(w, e, Connection{
		Target:  clientAddress,
		Channel: channelsRes.UdpGob,
	})
	ecs.AddComponent(w, e, Outbox[Message]{
		Data: make([]Message, 0),
	})
	ecs.AddComponent(w, e, Inbox[Message]{
		Data: make([]Message, 0),
	})

	t.Log("[Action] Entity created. Waiting for SessionMap update...")
	sch.RunUpdate(w) // next frame: update session map

	// Send Payload
	buf.Reset()
	encoder = gob.NewEncoder(&buf)
	_ = encoder.Encode(Message{Content: "Payload Data"})
	t.Log("[Action] Sending 'Payload Data'...")
	client.Write(buf.Bytes())

	time.Sleep(20 * time.Millisecond)
	sch.RunUpdate(w) // next frame: from cache to inbox

	// Get data
	inbox, ok := ecs.GetComponent[Inbox[Message]](w, e)
	if !ok {
		t.Fatal("Inbox component missing")
	}

	t.Logf("[Check] Inbox contains %d messages", len(inbox.Data))
	if len(inbox.Data) == 0 {
		// Debug the cache again
		c, _ := ecs.GetResource[cache](w)
		t.Logf("[Debug] Cache has %d packets. SessionLookup Size: %d", len(c.packets), len(sess.lookup))
		for k, v := range sess.lookup {
			t.Logf("   -> Map: %s = Entity %d", k, v)
		}
		t.Fatal("Inbox is empty, expected 1 message")
	}

	if inbox.Data[0].Content != "Payload Data" {
		t.Errorf("Expected 'Payload Data', got '%s'", inbox.Data[0].Content)
	}

	// Clean inbox
	inbox.Data = inbox.Data[:0]

	// Add message to outbox
	outbox, _ := ecs.GetComponent[Outbox[Message]](w, e)
	outbox.Data = append(outbox.Data, Message{Content: "Reply from ECS"})

	t.Log("[Action] Sending Reply...")
	sch.RunUpdate(w) // next frame: outbox to transport

	// Read from udp
	client.SetReadDeadline(time.Now().Add(time.Second * 1))
	readBuf := make([]byte, 1024)
	n, _, err := client.ReadFromUDP(readBuf)
	if err != nil {
		t.Fatalf("Client failed to receive data: %v", err)
	}

	// Decode response
	var receivedMsg Message
	decoder := gob.NewDecoder(bytes.NewReader(readBuf[:n]))
	if err := decoder.Decode(&receivedMsg); err != nil {
		t.Fatalf("Client failed to decode response: %v", err)
	}

	if receivedMsg.Content != "Reply from ECS" {
		t.Errorf("Client received wrong message: %s", receivedMsg.Content)
	}

	t.Log("Success! Round trip completed.")
}
