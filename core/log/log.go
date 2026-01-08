package log

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/vistormu/go-dsa/queue"

	"github.com/vistormu/xpeto/core/clock"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

type logger struct {
	sinks []Sink

	records *queue.RingQueue[Record]
	dropped uint64

	scratch []Record

	framesSinceFlush  int
	recordsSinceFlush int
}

func newLogger(s LoggerSettings) func() logger {
	return func() logger {
		if s.MaxRecords == 0 {
			s.MaxRecords = newLoggerSettings().MaxRecords
		}

		capacity := max(0, s.MaxRecords)

		q := queue.NewRingQueue[Record](max(1, capacity))

		return logger{
			records: q,
			scratch: make([]Record, 0, 1024),
			sinks:   make([]Sink, 0),
		}
	}
}

func log(w *ecs.World, level Level, msg string, fields []field, callerSkip int) {
	s := ecs.EnsureResource(w, newLoggerSettings)
	if level < s.MinLevel {
		return
	}

	if s.SilencedLevels.Contains(level) {
		return
	}

	l := ecs.EnsureResource(w, newLogger(*s))

	var frame uint64
	rClk, ok := ecs.GetResource[clock.RealClock](w)
	if ok {
		frame = rClk.Frame
	}

	var elapsed time.Duration
	vClk, ok := ecs.GetResource[clock.VirtualClock](w)
	if ok {
		elapsed = vClk.Elapsed
	}

	var id uint64
	label := "unnamed system"
	rs, ok := ecs.GetResource[schedule.RunningSystem](w)
	if ok {
		id = rs.Id
		label = rs.Label
	}

	var c Caller
	if s.CaptureCaller {
		skip := 2 + s.CallerSkip + callerSkip
		c = caller(skip)
	}

	r := Record{
		Level:       level,
		SystemId:    id,
		SystemLabel: label,
		Frame:       frame,
		Time:        elapsed,
		Caller:      c,
		Message:     msg,
		Fields:      fields,
	}

	if s.MaxRecords <= 0 {
		l.dropped++
		return
	}

	ok = l.records.Enqueue(r)
	if !ok {
		_, _ = l.records.Dequeue()
		_ = l.records.Enqueue(r)
		l.dropped++
	}
}

func flush(w *ecs.World) {
	s := ecs.EnsureResource(w, newLoggerSettings)
	l := ecs.EnsureResource(w, newLogger(*s))

	if l.records.Len() == 0 {
		return
	}

	l.scratch = l.scratch[:0]
	for r := range l.records.Drain() {
		l.scratch = append(l.scratch, r)
	}

	if len(l.scratch) == 0 {
		return
	}

	if l.dropped != 0 {
		lastFrame := l.scratch[len(l.scratch)-1].Frame
		l.scratch = append(l.scratch, Record{
			Level:       Warning,
			SystemId:    0,
			SystemLabel: "log",
			Frame:       lastFrame,
			Time:        0,
			Message:     "log buffer dropped records",
			Fields: []field{
				F("dropped", l.dropped),
				F("maxRecords", s.MaxRecords),
			},
		})
		l.dropped = 0
	}

	start := 0
	for start < len(l.scratch) {
		frame := l.scratch[start].Frame

		end := start + 1
		for end < len(l.scratch) && l.scratch[end].Frame == frame {
			end++
		}

		frameRecords := l.scratch[start:end]

		for _, s := range l.sinks {
			if s == nil {
				continue
			}
			s.Write(frame, frameRecords)
		}

		start = end
	}

	l.recordsSinceFlush += len(l.scratch)
	l.framesSinceFlush++

	// policy
	doFlush := false

	switch s.FlushMode {
	case FlushManual:
		return

	case FlushEveryFrame:
		doFlush = true

	case FlushEveryNFrames:
		n := s.FlushEveryNFrames
		if n > 0 && l.framesSinceFlush >= n {
			doFlush = true
		}

	case FlushEveryNRecords:
		n := s.FlushEveryNRecords
		if n > 0 && l.recordsSinceFlush >= n {
			doFlush = true
		}
	}

	// flush
	if !doFlush {
		return
	}

	for _, sink := range l.sinks {
		_ = sink.Flush()
		if s.SyncOnFlush {
			_ = sink.Sync()
		}
	}

	l.framesSinceFlush = 0
	l.recordsSinceFlush = 0
}

// ===
// API
// ===
func LogDebug(w *ecs.World, msg string, fields ...field) {
	log(w, Debug, msg, fields, 0)
}

func LogInfo(w *ecs.World, msg string, fields ...field) {
	log(w, Info, msg, fields, 0)
}

func LogWarning(w *ecs.World, msg string, fields ...field) {
	log(w, Warning, msg, fields, 0)
}

func LogError(w *ecs.World, msg string, fields ...field) {
	log(w, Error, msg, fields, 0)
}

func LogFatal(w *ecs.World, msg string, fields ...field) {
	log(w, Fatal, msg, fields, 0)
}

func AddSink(w *ecs.World, s Sink) {
	l, ok := ecs.GetResource[logger](w)
	if !ok || s == nil {
		return
	}
	l.sinks = append(l.sinks, s)
}

func ClearSinks(w *ecs.World) {
	l, ok := ecs.GetResource[logger](w)
	if !ok {
		return
	}
	l.sinks = l.sinks[:0]
}

func RecoverLog(w *ecs.World) func() {
	return func() {
		r := recover()
		if r == nil {
			return
		}

		LogFatal(w, "panic recovered",
			F("panic", fmt.Sprint(r)),
			F("stack", string(debug.Stack())),
		)
		flush(w)

		panic(r)
	}
}

func SetLoggerMaxRecords(w *ecs.World, n int) {
	s, ok := ecs.GetResource[LoggerSettings](w)
	if !ok {
		return
	}
	l, ok := ecs.GetResource[logger](w)
	if !ok {
		return
	}

	s.MaxRecords = n
	if n <= 0 {
		return
	}

	buf := make([]Record, 0, l.records.Len())
	for r := range l.records.Drain() {
		buf = append(buf, r)
	}
	if len(buf) > n {
		buf = buf[len(buf)-n:]
	}

	l.records = queue.NewRingQueue[Record](n)
	for _, r := range buf {
		_ = l.records.Enqueue(r)
	}
}
