package log

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/vistormu/xpeto/core/clock"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

type sinkCall struct {
	frame   uint64
	records []Record
}

type captureSink struct {
	writeCalls []sinkCall
	flushCalls int
	syncCalls  int
}

func (s *captureSink) Write(frame uint64, records []Record) {
	cp := make([]Record, len(records))
	copy(cp, records)
	s.writeCalls = append(s.writeCalls, sinkCall{frame: frame, records: cp})
}

func (s *captureSink) Flush() error {
	s.flushCalls++
	return nil
}

func (s *captureSink) Sync() error {
	s.syncCalls++
	return nil
}

type errorSink struct {
	writeCalls int
	flushCalls int
	syncCalls  int
}

func (s *errorSink) Write(frame uint64, records []Record) { s.writeCalls++ }

func (s *errorSink) Flush() error {
	s.flushCalls++
	return fmt.Errorf("flush error")
}

func (s *errorSink) Sync() error {
	s.syncCalls++
	return fmt.Errorf("sync error")
}

type memJSONSink struct {
	buf bytes.Buffer
}

func (s *memJSONSink) Write(frame uint64, records []Record) {
	for _, r := range records {
		// minimal ndjson like output, for test purposes
		fmt.Fprintf(&s.buf, `{"frame":%d,"level":"%s","msg":%q}`+"\n", frame, r.Level.String(), r.Message)
	}
}

func (s *memJSONSink) Flush() error { return nil }

func (s *memJSONSink) Sync() error { return nil }

func newTestWorld() *ecs.World {
	w := ecs.NewWorld()
	sch := schedule.NewScheduler()

	ecs.AddResource(w, schedule.RunningSystem{})
	clock.Pkg(w, sch)
	Pkg(w, sch)

	return w
}

func setFrame(w *ecs.World, frame uint64) {
	r, ok := ecs.GetResource[clock.RealClock](w)
	if ok && r != nil {
		r.Frame = frame
	}
}

func TestLog_MinLevelAndSilencedLevels(t *testing.T) {
	w := newTestWorld()
	s := &captureSink{}
	AddSink(w, s)

	SetLogLevel(w, Warning)
	LogInfo(w, "ignored")
	LogWarning(w, "kept")
	flush(w)

	if len(s.writeCalls) != 1 {
		t.Fatalf("expected 1 Write call, got %d", len(s.writeCalls))
	}
	if got := len(s.writeCalls[0].records); got != 1 {
		t.Fatalf("expected 1 record, got %d", got)
	}
	if s.writeCalls[0].records[0].Message != "kept" {
		t.Fatalf("unexpected message: %q", s.writeCalls[0].records[0].Message)
	}

	// silence Warning, allow Error
	cfg, _ := ecs.GetResource[LoggerSettings](w)
	cfg.MinLevel = Debug
	cfg.SilencedLevels.Add(Warning)

	LogWarning(w, "silenced")
	LogError(w, "not silenced")
	flush(w)

	// one more write call for frame 0
	if len(s.writeCalls) != 2 {
		t.Fatalf("expected 2 Write calls, got %d", len(s.writeCalls))
	}
	recs := s.writeCalls[1].records
	if len(recs) != 1 || recs[0].Message != "not silenced" {
		t.Fatalf("expected only the Error record, got %+v", recs)
	}
}

func TestFlush_GroupsByFrame(t *testing.T) {
	w := newTestWorld()
	s := &captureSink{}
	AddSink(w, s)

	setFrame(w, 10)
	LogInfo(w, "a")
	setFrame(w, 11)
	LogInfo(w, "b")
	setFrame(w, 11)
	LogInfo(w, "c")

	flush(w)

	if len(s.writeCalls) != 2 {
		t.Fatalf("expected 2 frame groups, got %d", len(s.writeCalls))
	}
	if s.writeCalls[0].frame != 10 || s.writeCalls[1].frame != 11 {
		t.Fatalf("unexpected frames: %d, %d", s.writeCalls[0].frame, s.writeCalls[1].frame)
	}
	if got := len(s.writeCalls[0].records); got != 1 {
		t.Fatalf("expected 1 record in frame 10, got %d", got)
	}
	if got := len(s.writeCalls[1].records); got != 2 {
		t.Fatalf("expected 2 records in frame 11, got %d", got)
	}
}

func TestFlushPolicy_EveryFrame(t *testing.T) {
	w := newTestWorld()
	s := &captureSink{}
	AddSink(w, s)

	cfg, _ := ecs.GetResource[LoggerSettings](w)
	cfg.FlushMode = FlushEveryFrame
	cfg.SyncOnFlush = false

	setFrame(w, 1)
	LogInfo(w, "a")
	flush(w)

	if s.flushCalls != 1 {
		t.Fatalf("expected 1 Flush call, got %d", s.flushCalls)
	}
	if s.syncCalls != 0 {
		t.Fatalf("expected 0 Sync calls, got %d", s.syncCalls)
	}
}

func TestFlushPolicy_EveryNFrames(t *testing.T) {
	w := newTestWorld()
	s := &captureSink{}
	AddSink(w, s)

	cfg, _ := ecs.GetResource[LoggerSettings](w)
	cfg.FlushMode = FlushEveryNFrames
	cfg.FlushEveryNFrames = 2
	cfg.SyncOnFlush = true

	setFrame(w, 1)
	LogInfo(w, "a")
	flush(w)

	if s.flushCalls != 0 || s.syncCalls != 0 {
		t.Fatalf("expected no flush yet, got flush=%d sync=%d", s.flushCalls, s.syncCalls)
	}

	setFrame(w, 2)
	LogInfo(w, "b")
	flush(w)

	if s.flushCalls != 1 || s.syncCalls != 1 {
		t.Fatalf("expected flush+sync after 2 frames, got flush=%d sync=%d", s.flushCalls, s.syncCalls)
	}
}

func TestFlushPolicy_EveryNRecords(t *testing.T) {
	w := newTestWorld()
	s := &captureSink{}
	AddSink(w, s)

	cfg, _ := ecs.GetResource[LoggerSettings](w)
	cfg.FlushMode = FlushEveryNRecords
	cfg.FlushEveryNRecords = 3
	cfg.SyncOnFlush = false

	setFrame(w, 1)
	LogInfo(w, "a")
	LogInfo(w, "b")
	flush(w)

	if s.flushCalls != 0 {
		t.Fatalf("expected no flush yet, got %d", s.flushCalls)
	}

	setFrame(w, 2)
	LogInfo(w, "c")
	flush(w)

	if s.flushCalls != 1 {
		t.Fatalf("expected flush after 3 records, got %d", s.flushCalls)
	}
}

func TestSetLoggerMaxRecords_KeepsMostRecentAndAddsDroppedWarning(t *testing.T) {
	w := newTestWorld()
	s := &captureSink{}
	AddSink(w, s)

	SetLoggerMaxRecords(w, 2)

	setFrame(w, 1)
	LogInfo(w, "a")
	LogInfo(w, "b")
	LogInfo(w, "c")

	flush(w)

	if len(s.writeCalls) != 1 {
		t.Fatalf("expected 1 Write call, got %d", len(s.writeCalls))
	}

	var msgs []string
	for _, r := range s.writeCalls[0].records {
		msgs = append(msgs, r.Message)
	}

	// "a" should be dropped, "b" and "c" kept, plus warning record.
	hasB := false
	hasC := false
	hasDropped := false
	for _, m := range msgs {
		if m == "b" {
			hasB = true
		}
		if m == "c" {
			hasC = true
		}
		if m == "log buffer dropped records" {
			hasDropped = true
		}
	}

	if !hasB || !hasC || !hasDropped {
		t.Fatalf("unexpected messages: %v", msgs)
	}
}

func TestRecoverLog_LogsAndRepanics(t *testing.T) {
	w := newTestWorld()
	s := &captureSink{}
	AddSink(w, s)

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic to propagate")
			}
		}()
		defer RecoverLog(w)()
		panic("boom")
	}()

	found := false
	for _, c := range s.writeCalls {
		for _, r := range c.records {
			if r.Level == Fatal && r.Message == "panic recovered" {
				found = true
				break
			}
		}
	}
	if !found {
		t.Fatalf("expected a fatal panic recovered record")
	}
}

func TestFlush_IgnoresSinkErrors(t *testing.T) {
	w := newTestWorld()
	s := &errorSink{}
	AddSink(w, s)

	cfg, _ := ecs.GetResource[LoggerSettings](w)
	cfg.FlushMode = FlushEveryFrame
	cfg.SyncOnFlush = true

	setFrame(w, 1)
	LogInfo(w, "a")
	flush(w)

	if s.flushCalls != 1 || s.syncCalls != 1 {
		t.Fatalf("expected sink flush+sync calls, got flush=%d sync=%d", s.flushCalls, s.syncCalls)
	}
}

func TestSink_NDJSONShape(t *testing.T) {
	w := newTestWorld()
	j := &memJSONSink{}
	AddSink(w, j)

	cfg, _ := ecs.GetResource[LoggerSettings](w)
	cfg.FlushMode = FlushEveryFrame

	setFrame(w, 7)
	LogInfo(w, "hello")
	LogError(w, "fail")
	flush(w)

	out := j.buf.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %q", len(lines), out)
	}
	if !strings.Contains(lines[0], `"frame":7`) || !strings.Contains(lines[0], `"msg":"hello"`) {
		t.Fatalf("unexpected line 1: %q", lines[0])
	}
	if !strings.Contains(lines[1], `"frame":7`) || !strings.Contains(lines[1], `"msg":"fail"`) {
		t.Fatalf("unexpected line 2: %q", lines[1])
	}
}
