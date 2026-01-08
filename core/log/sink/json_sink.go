package sink

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/vistormu/xpeto/core/log"
)

type JsonSink struct {
	mu sync.Mutex

	out io.Writer
	b   *bufio.Writer
}

func NewJSONSink(out io.Writer) *JsonSink {
	if out == nil {
		out = os.Stdout
	}
	return &JsonSink{
		out: out,
		b:   bufio.NewWriterSize(out, 64*1024),
	}
}

func (s *JsonSink) Write(frame uint64, records []log.Record) {
	if len(records) == 0 {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.b == nil {
		return
	}

	for _, r := range records {
		obj := map[string]any{
			"frame":    frame,
			"level":    r.Level.String(),
			"msg":      r.Message,
			"system":   r.SystemLabel,
			"systemId": r.SystemId,
			"time":     r.Time.String(),
		}

		if r.Caller.File != "" {
			obj["caller"] = map[string]any{
				"file": r.Caller.File,
				"line": r.Caller.Line,
				"func": r.Caller.Func,
			}
		}

		if len(r.Fields) != 0 {
			fields := make(map[string]any, len(r.Fields))
			for _, f := range r.Fields {
				fields[f.Key()] = f.Value()
			}
			obj["fields"] = fields
		}

		b, err := json.Marshal(obj)
		if err != nil {
			fmt.Fprintf(s.b, "{\"frame\":%d,\"level\":\"error\",\"msg\":%q}\n", frame, err.Error())
			continue
		}

		_, _ = s.b.Write(b)
		_, _ = s.b.WriteString("\n")
	}
}

func (s *JsonSink) Flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.b == nil {
		return nil
	}
	return s.b.Flush()
}

func (s *JsonSink) Sync() error {
	return s.Flush()
}
