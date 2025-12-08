package net

import (
	"bytes"
	"runtime"
	"sync"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/log"
	"github.com/vistormu/xpeto/pkg/net/codec"
	"github.com/vistormu/xpeto/pkg/net/transport"
)

// =====
// types
// =====
type Inbox[T any] struct {
	Data []T
}

type cache struct {
	packets      []transport.Packet
	packetCodecs []codec.Codec
}

// =======
// helpers
// =======
type decodeResult[T any] struct {
	entity  ecs.Entity
	message T
	err     error
}

func receive[T any](w *ecs.World) {
	cache, ok := ecs.GetResource[cache](w)
	if !ok || len(cache.packets) == 0 {
		return
	}

	sessions, _ := ecs.GetResource[session](w)

	nWorkers := runtime.GOMAXPROCS(0)
	n := len(cache.packets)

	if n < 100 {
		nWorkers = 1
	}

	results := make(chan decodeResult[T], n)
	var wg sync.WaitGroup

	chunkSize := (n + nWorkers - 1) / nWorkers

	for i := range nWorkers {
		start := i * chunkSize
		end := start + chunkSize

		end = min(end, n)

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for j := start; j < end; j++ {
				packet := cache.packets[j]

				e, ok := sessions.lookup[packet.Sender]
				if !ok {
					continue
				}

				reader := bytes.NewReader(packet.Payload)
				var msg T
				cd := cache.packetCodecs[j]

				err := cd.Decode(reader, &msg)
				if err == nil {
					results <- decodeResult[T]{entity: e, message: msg}
				} else {
					results <- decodeResult[T]{entity: e, err: err}
				}
			}
		}(start, end)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		if r.err != nil {
			log.LogError(w, "error reading from stream", log.F("error", r.err))
			continue
		}

		inbox, ok := ecs.GetComponent[Inbox[T]](w, r.entity)
		if !ok {
			log.LogWarning(w, "a message arrived but the component net.Inbox[T] was missing", log.F("entity", r.entity))
			continue
		}
		if inbox.Data == nil {
			inbox.Data = make([]T, 0)
		}

		inbox.Data = append(inbox.Data, r.message)
	}
}
