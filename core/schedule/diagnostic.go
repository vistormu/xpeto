package schedule

import (
	"github.com/vistormu/go-dsa/queue"
)

// ==========
// diagnostic
// ==========
type Diagnostic struct {
	Message string
	Id      uint64
	Label   string
	Stage   stage
}

// ======
// buffer
// ======
type logger struct {
	diagnostics *queue.Queue[Diagnostic]
}

func newLogger() *logger {
	return &logger{
		diagnostics: queue.NewQueue[Diagnostic](),
	}
}

func (l *logger) add(message string, id uint64, label string, stage stage) {
	l.diagnostics.Enqueue(Diagnostic{
		Message: message,
		Id:      id,
		Label:   label,
		Stage:   stage,
	})
}
