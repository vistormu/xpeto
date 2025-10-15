package render

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/pkg/window"
)

type Renderable interface {
	SortKey() uint64
	Draw(screen *ebiten.Image)
}

func draw(w *ecs.World) {
	screen, _ := ecs.GetResource[window.Screen](w)
	e, _ := ecs.GetResource[extractor](w)

	// extraction
	for phase, fns := range e.extractors {
		buf := e.renderables[phase][:0]
		for _, fn := range fns {
			buf = append(buf, fn(w)...)
		}
		e.renderables[phase] = buf
	}

	// sort
	for phase := range e.renderables {
		items := e.renderables[phase]
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].SortKey() < items[j].SortKey()
		})
	}

	// render
	// TODO: phase order
	for phase := range e.renderables {
		for _, r := range e.renderables[phase] {
			r.Draw(screen)
		}
	}
}
