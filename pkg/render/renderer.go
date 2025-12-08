package render

import (
	"reflect"
	"sort"

	"github.com/vistormu/go-dsa/hashmap"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/log"
)

// =========
// functions
// =========
type extractionFn[T any] = func(*ecs.World) []T
type sortFn[T any] = func(v T) uint64
type renderFn[C, T any] = func(canvas *C, v T)
type renderFnWrapper[C any] func(r *renderer[C], w *ecs.World)

// ==========
// renderable
// ==========
type renderable struct {
	key     uint64
	batchId int
	index   int
}

// =====
// batch
// =====
type batch[C any] interface {
	reset()
	render(canvas *C, index int)
}

type concreteBatch[C any, T any] struct {
	data     []T
	renderFn renderFn[C, T]
}

func (b *concreteBatch[C, T]) reset() {
	b.data = b.data[:0]
}

func (b *concreteBatch[C, T]) render(canvas *C, index int) {
	b.renderFn(canvas, b.data[index])
}

// ========
// renderer
// ========
type renderer[C any] struct {
	extractionFns *hashmap.TypeMap
	sortFns       *hashmap.TypeMap
	renderFns     map[RenderStage][]renderFnWrapper[C]
	batches       []batch[C]
	renderables   []renderable
}

func newRenderer[C any]() *renderer[C] {
	return &renderer[C]{
		extractionFns: hashmap.NewTypeMap(),
		sortFns:       hashmap.NewTypeMap(),
		renderFns:     make(map[RenderStage][]renderFnWrapper[C]),
		batches:       make([]batch[C], 0),
		renderables:   make([]renderable, 0, 1000),
	}
}

// =======
// systems
// =======
func render[C any](w *ecs.World) {
	canvas, _ := ecs.GetResource[C](w)
	r, _ := ecs.GetResource[renderer[C]](w)

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
			r.batches[re.batchId].render(canvas, re.index)
		}
	}
}

// ===
// API
// ===
func AddExtractionFn[C, T any](w *ecs.World, fn extractionFn[T]) {
	r, ok := ecs.GetResource[renderer[C]](w)
	if !ok {
		log.LogError(w, renderPkgErr, log.F("function", "AddExtractionFunction[C, T any]"), log.F("canvas", reflect.TypeFor[C]().String()), log.F("type", reflect.TypeFor[T]().String()))
		return
	}

	hashmap.Add(r.extractionFns, fn)
}

func AddSortFn[C, T any](w *ecs.World, fn sortFn[T]) {
	r, ok := ecs.GetResource[renderer[C]](w)
	if !ok {
		log.LogError(w, renderPkgErr, log.F("function", "AddSort[C, T any]"), log.F("canvas", reflect.TypeFor[C]().String()), log.F("type", reflect.TypeFor[T]().String()))
		return
	}

	hashmap.Add(r.sortFns, fn)
}

func AddRenderFn[C, T any](w *ecs.World, stage RenderStage, fn renderFn[C, T]) {
	r, ok := ecs.GetResource[renderer[C]](w)
	if !ok {
		log.LogError(w, renderPkgErr, log.F("function", "AddRenderFn[C, T any]"), log.F("canvas", reflect.TypeFor[C]().String()), log.F("type", reflect.TypeFor[T]().String()))
		return
	}

	// create batch
	newBatch := &concreteBatch[C, T]{
		data:     make([]T, 0, 100),
		renderFn: fn,
	}
	batchId := len(r.batches)
	r.batches = append(r.batches, newBatch)

	exFn, okEx := hashmap.Get[extractionFn[T]](r.extractionFns)
	sortFn, okSort := hashmap.Get[sortFn[T]](r.sortFns)
	if !okEx {
		log.LogError(w, missingExFnErr, log.F("function", "AddRenderFn"), log.F("type", reflect.TypeFor[T]().String()))
		return
	}
	if !okSort {
		log.LogError(w, missingSortFnErr, log.F("function", "AddRenderFn"), log.F("type", reflect.TypeFor[T]().String()))
		return
	}

	_, ok = r.renderFns[stage]
	if !ok {
		r.renderFns[stage] = make([]renderFnWrapper[C], 0)
	}

	r.renderFns[stage] = append(r.renderFns[stage], func(r *renderer[C], w *ecs.World) {
		renderables := (*exFn)(w)
		if len(renderables) == 0 {
			return
		}

		newBatch.reset()
		currentBatchIndex := 0

		for _, re := range renderables {
			key := (*sortFn)(re)
			newBatch.data = append(newBatch.data, re)

			r.renderables = append(r.renderables, renderable{
				key:     key,
				batchId: batchId,
				index:   currentBatchIndex,
			})
			currentBatchIndex++
		}
	})
}
