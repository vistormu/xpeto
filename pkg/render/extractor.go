package render

type Extractor struct {
	extractors  map[Phase][]ExtractionFn
	renderables map[Phase][]Renderable
}

func NewExtractor() *Extractor {
	return &Extractor{
		extractors:  make(map[Phase][]ExtractionFn),
		renderables: make(map[Phase][]Renderable),
	}
}

func (e *Extractor) AddExtractionFn(phase Phase, fn ExtractionFn) {
	_, ok := e.extractors[phase]
	if !ok {
		e.extractors[phase] = []ExtractionFn{}
	}

	e.extractors[phase] = append(e.extractors[phase], fn)
}
