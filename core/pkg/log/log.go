package log

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/pkg/time"
)

type logger struct {
	minLevel Level
	records  []record
	sinks    []sink
}

func newLogger() logger {
	return logger{
		minLevel: Debug,
		records:  make([]record, 0),
		sinks:    make([]sink, 0),
	}
}

func log(w *ecs.World, level Level, msg string, fields []field) {
	l, _ := ecs.GetResource[logger](w)
	if level < l.minLevel {
		return
	}

	rClk, _ := ecs.GetResource[time.RealClock](w)
	vClk, _ := ecs.GetResource[time.VirtualClock](w)

	id := ecs.GetSystemId(w)

	r := record{
		level:    level,
		systemId: id,
		frame:    vClk.Frame,
		time:     rClk.Elapsed,
		message:  msg,
		fields:   fields,
	}

	l.records = append(l.records, r)
}

func flush(w *ecs.World) {
	l, _ := ecs.GetResource[logger](w)

	if len(l.records) == 0 {
		return
	}

	start := 0
	for start < len(l.records) {
		frame := l.records[start].frame

		end := start + 1
		for end < len(l.records) && l.records[end].frame == frame {
			end++
		}

		frameRecords := l.records[start:end]

		for _, s := range l.sinks {
			s.write(frame, frameRecords)
		}

		start = end
	}

	l.records = l.records[:0]
}

// ===
// API
// ===
func SetLogLevel(w *ecs.World, l Level) {
	logger, _ := ecs.GetResource[logger](w)
	logger.minLevel = l
}

func LogDebug(w *ecs.World, msg string, fields ...field) {
	log(w, Debug, msg, fields)
}

func LogInfo(w *ecs.World, msg string, fields ...field) {
	log(w, Info, msg, fields)
}

func LogWarning(w *ecs.World, msg string, fields ...field) {
	log(w, Warning, msg, fields)
}

func LogError(w *ecs.World, msg string, fields ...field) {
	log(w, Error, msg, fields)
}

func LogFatal(w *ecs.World, msg string, fields ...field) {
	log(w, Fatal, msg, fields)
}
