package sink

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/vistormu/xpeto/core/log"
)

type FileSink struct {
	mu sync.Mutex

	path string
	f    *os.File
	w    *bufio.Writer
}

func NewFileSink(path string) (*FileSink, error) {
	if path == "" {
		return nil, fmt.Errorf("file sink: empty path")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}

	return &FileSink{
		path: path,
		f:    f,
		w:    bufio.NewWriterSize(f, 64*1024),
	}, nil
}

func (s *FileSink) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.f == nil {
		return nil
	}

	_ = s.w.Flush()
	err := s.f.Close()
	s.f = nil
	s.w = nil
	return err
}

func (s *FileSink) Write(frame uint64, records []log.Record) {
	if len(records) == 0 {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.f == nil || s.w == nil {
		return
	}

	for _, r := range records {
		if r.Caller.File != "" {
			fmt.Fprintf(s.w,
				"frame=%d level=%s msg=%q system=%q id=%d time=%s caller=%s:%d func=%q",
				frame,
				r.Level.String(),
				r.Message,
				r.SystemLabel,
				r.SystemId,
				r.Time,
				r.Caller.File,
				r.Caller.Line,
				r.Caller.Func,
			)
		} else {
			fmt.Fprintf(s.w,
				"frame=%d level=%s msg=%q system=%q id=%d time=%s",
				frame,
				r.Level.String(),
				r.Message,
				r.SystemLabel,
				r.SystemId,
				r.Time,
			)
		}

		for _, f := range r.Fields {
			fmt.Fprintf(s.w, " %s=%v", f.Key(), f.Value())
		}
		_, _ = s.w.WriteString("\n")
	}
}

func (s *FileSink) Flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.w == nil {
		return nil
	}
	return s.w.Flush()
}

func (s *FileSink) Sync() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.f == nil || s.w == nil {
		return nil
	}
	if err := s.w.Flush(); err != nil {
		return err
	}
	return s.f.Sync()
}
