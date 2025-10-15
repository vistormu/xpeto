package render

import (
	"github.com/vistormu/xpeto/core/ecs"
)

type extractionFn = func(*ecs.World) []Renderable

type Phase int

const (
	Transparent Phase = iota
	Opaque
	UI
	PostFx
)

type extractor struct {
	extractors  map[Phase][]extractionFn
	renderables map[Phase][]Renderable
}

func newExtractor() *extractor {
	return &extractor{
		extractors:  make(map[Phase][]extractionFn),
		renderables: make(map[Phase][]Renderable),
	}
}

// ===
// API
// ===
func AddExtractionFn(w *ecs.World, phase Phase, fn extractionFn) {
	e, _ := ecs.GetResource[extractor](w)

	_, ok := e.extractors[phase]
	if !ok {
		e.extractors[phase] = []extractionFn{}
	}

	e.extractors[phase] = append(e.extractors[phase], fn)
}
