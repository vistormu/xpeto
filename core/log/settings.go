package log

import (
	"github.com/vistormu/go-dsa/set"

	"github.com/vistormu/xpeto/core/ecs"
)

// ====
// mode
// ====
type FlushMode uint8

const (
	FlushManual FlushMode = iota
	FlushEveryFrame
	FlushEveryNFrames
	FlushEveryNRecords
)

// ========
// settings
// ========
type LoggerSettings struct {
	MinLevel       Level
	SilencedLevels *set.HashSet[Level]
	MaxRecords     int

	CaptureCaller bool
	CallerSkip    int

	FlushMode          FlushMode
	FlushEveryNFrames  int
	FlushEveryNRecords int
	SyncOnFlush        bool
}

func newLoggerSettings() LoggerSettings {
	return LoggerSettings{
		MinLevel:       Debug,
		SilencedLevels: set.NewHashSet[Level](),
		MaxRecords:     16_384,
		CaptureCaller:  false,
		CallerSkip:     0,
		FlushMode:      FlushEveryFrame,
	}
}

// ===
// API
// ===
func SetLogLevel(w *ecs.World, l Level) {
	s := ecs.EnsureResource(w, newLoggerSettings)
	s.MinLevel = l
}

func FlushLoggerManually(w *ecs.World) {
	s := ecs.EnsureResource(w, newLoggerSettings)
	s.FlushMode = FlushManual
}

func FlushLoggerEveryFrame(w *ecs.World) {
	s := ecs.EnsureResource(w, newLoggerSettings)
	s.FlushMode = FlushEveryFrame
}

func FlushLoggerEveryNFrames(w *ecs.World, n int) {
	s := ecs.EnsureResource(w, newLoggerSettings)
	s.FlushMode = FlushEveryNFrames
	s.FlushEveryNFrames = n
}

func FlushLoggerEveryNRecords(w *ecs.World, n int) {
	s := ecs.EnsureResource(w, newLoggerSettings)
	s.FlushMode = FlushEveryNRecords
	s.FlushEveryNRecords = n
}

func SetSyncOnFlush(w *ecs.World, enabled bool) {
	s := ecs.EnsureResource(w, newLoggerSettings)
	s.SyncOnFlush = enabled
}

func SilenceLevels(w *ecs.World, levels ...Level) {
	s := ecs.EnsureResource(w, newLoggerSettings)
	for _, l := range levels {
		s.SilencedLevels.Add(l)
	}
}

func UnsilenceLevels(w *ecs.World, levels ...Level) {
	s := ecs.EnsureResource(w, newLoggerSettings)
	for _, l := range levels {
		s.SilencedLevels.Remove(l)
	}
}
