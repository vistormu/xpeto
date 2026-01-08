package asset

import (
	"fmt"
	"io/fs"
	"runtime"
	"sync"

	"github.com/vistormu/go-dsa/queue"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/log"
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
	requests *queue.Queue[request]
	pending  []result
	sem      chan struct{}
}

func newLoader() loader {
	n := max(runtime.GOMAXPROCS(0)*2, 4)

	return loader{
		requests: queue.NewQueue[request](),
		pending:  make([]result, 0),
		sem:      make(chan struct{}, n),
	}
}

func (l *loader) acquire() {
	l.sem <- struct{}{}
}

func (l *loader) release() {
	<-l.sem
}

func (l *loader) enqueue(req request) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.requests.Enqueue(req)
}

func (l *loader) drainRequests() []request {
	l.mu.Lock()
	defer l.mu.Unlock()

	out := make([]request, 0, l.requests.Len())
	for req := range l.requests.Drain() {
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
	if base == "" || rel == "" {
		return nil, fmt.Errorf("invalid asset path %q (expected base/rel)", req.path)
	}

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
		if !s.population.has(req.asset) {
			continue
		}

		s.population.setLoading(req.asset)

		l.acquire()
		go func(r request) {
			defer l.release()

			data, err := readFile(s, r)
			l.addResult(result{
				req:  r,
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
		if !s.population.has(res.req.asset) {
			continue
		}

		if res.err != nil {
			s.population.setFailed(res.req.asset, res.err)

			log.LogError(w, "failed loading asset",
				log.F("path", res.req.path),
				log.F("error", res.err.Error()),
			)
			continue
		}

		_, _, ext := splitPath(res.req.path)
		fn, ok := s.loaders[ext]
		if !ok {
			err := fmt.Errorf("loader not found for the extension %s", ext)
			s.population.setFailed(res.req.asset, err)

			log.LogError(w, "loader not found for extension. this error shouldn't have happened", log.F("ext", ext))
			continue
		}

		err := fn(s, res.req, res.data)
		if err != nil {
			s.population.setFailed(res.req.asset, err)
			log.LogError(w, "the asset loader returned an error",
				log.F("path", res.req.path),
				log.F("error", err.Error()),
			)
			continue
		}

		s.population.setLoaded(res.req.asset)
	}
}
