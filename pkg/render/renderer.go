package render

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/go-dsa/hashmap"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/pkg/log"
)

// =====
// types
// =====
type ExtractionFn[T any] = func(*ecs.World) []T
type SortFn[T any] = func(v T) uint64
type RenderFn[T any] = func(screen *ebiten.Image, v T)

type renderable struct {
	key  uint64
	draw func(*ebiten.Image)
}

type renderFn func(r *renderer, w *ecs.World)

// ========
// renderer
// ========
type renderer struct {
	extractionFns *hashmap.TypeMap
	sortFns       *hashmap.TypeMap
	renderFns     map[RenderStage][]renderFn
	renderables   []renderable
}

func newRenderer() renderer {
	return renderer{
		extractionFns: hashmap.NewTypeMap(),
		sortFns:       hashmap.NewTypeMap(),
		renderFns:     make(map[RenderStage][]renderFn),
		renderables:   make([]renderable, 0),
	}
}

// =======
// systems
// =======
func draw(w *ecs.World) {
	screen, _ := ecs.GetResource[ebiten.Image](w)
	r, _ := ecs.GetResource[renderer](w)

	for _, stage := range stagesOrder {
		fns, ok := r.renderFns[stage]
		if !ok || len(fns) == 0 {
			continue
		}

		r.renderables = r.renderables[:0]

		for _, fn := range fns {
			fn(r, w)
		}

		sort.SliceStable(r.renderables, func(i, j int) bool {
			return r.renderables[i].key < r.renderables[j].key
		})

		for _, re := range r.renderables {
			re.draw(screen)
		}
	}
}

// ===
// API
// ===
func AddExtractionFn[T any](w *ecs.World, fn ExtractionFn[T]) {
	r, ok := ecs.GetResource[renderer](w)
	if !ok {
		log.LogError(w, "could not execute AddExtractionFn: render.Pkg not added")
		return
	}

	hashmap.Add(r.extractionFns, fn)
}

func AddSortFn[T any](w *ecs.World, fn SortFn[T]) {
	r, ok := ecs.GetResource[renderer](w)
	if !ok {
		log.LogError(w, "could not execute AddSortFn: render.Pkg not added")
		return
	}

	hashmap.Add(r.sortFns, fn)
}

func AddRenderFn[T any](w *ecs.World, stage RenderStage, fn RenderFn[T]) {
	r, ok := ecs.GetResource[renderer](w)
	if !ok {
		log.LogError(w, "cannot execute AddRenderFn: render.Pkg not included")
		return
	}

	_, ok = r.renderFns[stage]
	if !ok {
		r.renderFns[stage] = make([]renderFn, 0)
	}

	exFn, ok := hashmap.Get[ExtractionFn[T]](r.extractionFns)
	if !ok {
		log.LogError(w, "AddExtractionFn should be executed before AddRenderFn")
		return
	}

	sortFn, ok := hashmap.Get[SortFn[T]](r.sortFns)
	if !ok {
		log.LogError(w, "AddSortFn should be executed before AddRenderFn")
		return
	}

	r.renderFns[stage] = append(r.renderFns[stage], func(r *renderer, w *ecs.World) {
		renderables := (*exFn)(w)
		if len(renderables) == 0 {
			return
		}

		for _, re := range renderables {
			reCopy := re
			key := (*sortFn)(reCopy)

			r.renderables = append(r.renderables, renderable{
				key: key,
				draw: func(screen *ebiten.Image) {
					fn(screen, reCopy)
				},
			})
		}
	})
}
