package asset

import (
	"fmt"
	"io/fs"
	"sync"

	"github.com/vistormu/go-dsa/queue"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/pkg/log"
)

// =====
// types
// =====
type LoaderFn[T any] = func([]byte, string) (*T, error)

type loaderHandler = func(*server, request, []byte) error

type request struct {
	path  string
	asset Asset
}

type result struct {
	req  request
	data []byte
	err  error
}

// ======
// loader
// ======
type loader struct {
	mu       sync.Mutex
	requests *queue.QueueArray[request]
	pending  []result
}

func newLoader() loader {
	return loader{
		requests: queue.NewQueueArray[request](),
		pending:  make([]result, 0),
	}
}

func (l *loader) enqueue(req request) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.requests.Enqueue(req)
}

func (l *loader) drainRequests() []request {
	l.mu.Lock()
	defer l.mu.Unlock()

	out := make([]request, 0, l.requests.Length())
	for !l.requests.IsEmpty() {
		req, _ := l.requests.Dequeue()
		out = append(out, req)
	}

	return out
}

func (l *loader) addResult(res result) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.pending = append(l.pending, res)
}

func (l *loader) drainResults() []result {
	l.mu.Lock()
	defer l.mu.Unlock()

	out := l.pending
	l.pending = nil

	return out
}

// =======
// helpers
// =======
func readFile(s *server, req request) ([]byte, error) {
	base, rel, _ := splitPath(req.path)

	fsys, ok := s.staticFS[base]
	if !ok {
		return nil, fmt.Errorf("could not find the filesystem for base %s. this error shouldn't have hapenned", base)
	}

	data, err := fs.ReadFile(fsys, rel)
	if err != nil {
		return nil, fmt.Errorf("cannot read asset %q (%s:%s): %w", req.path, base, rel, err)
	}

	return data, nil
}

// =======
// systems
// =======
func readRequests(w *ecs.World) {
	s, _ := ecs.GetResource[server](w)
	l, _ := ecs.GetResource[loader](w)

	for _, req := range l.drainRequests() {
		go func(r request) {
			data, err := readFile(s, req)
			l.addResult(result{
				req:  req,
				data: data,
				err:  err,
			})
		}(req)
	}
}

func loadResults(w *ecs.World) {
	s, _ := ecs.GetResource[server](w)
	l, _ := ecs.GetResource[loader](w)

	for _, res := range l.drainResults() {
		if res.err != nil {
			log.LogError(w, "failed loading asset",
				log.F("path", res.req.path),
				log.F("error", res.err.Error()),
			)
			continue
		}

		_, _, ext := splitPath(res.req.path)
		fn, ok := s.loaders[ext]
		if !ok {
			log.LogError(w, "loader not found for extension. this error shouldn't have happened", log.F("ext", ext))
			continue
		}

		err := fn(s, res.req, res.data)
		if err != nil {
			log.LogError(w, "the asset loader returned an error",
				log.F("path", res.req.path),
				log.F("error", err.Error()),
			)
		}
	}
}
