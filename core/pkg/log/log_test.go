package log

import (
	"testing"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/pkg/time"
)

// simple sink that stores what it receives
type testSink struct {
	frames []uint64
	logs   [][]record
}

func (s *testSink) write(frame uint64, records []record) {
	s.frames = append(s.frames, frame)

	cp := make([]record, len(records))
	copy(cp, records)

	s.logs = append(s.logs, cp)
}

func dummySystem(w *ecs.World) {
	LogWarning(w, "test-message", F("hello", "world"))
}

// test
func TestLoggerFlush(t *testing.T) {
	w := ecs.NewWorld()
	ecs.AddResource(w, time.RealClock{})
	ecs.AddResource(w, time.VirtualClock{
		Frame: 42,
	})

	s := &testSink{}

	l := logger{
		minLevel: Warning,
		sinks:    []sink{s, &debugSink{}},
	}
	ecs.AddResource(w, l)

	dummySystem(w)

	flush(w)

	// assertions
	if len(s.frames) != 1 {
		t.Fatalf("expected 1 frame, got %d", len(s.frames))
	}

	if s.frames[0] != 42 {
		t.Fatalf("expected frame 42, got %d", s.frames[0])
	}

	if len(s.logs) != 1 || len(s.logs[0]) != 1 {
		t.Fatalf("expected 1 log entry, got %v", s.logs)
	}

	r := s.logs[0][0]

	if r.message != "test-message" {
		t.Fatalf("wrong message: %s", r.message)
	}
}
