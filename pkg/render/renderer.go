package render

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vistormu/xpeto/internal/core"
)

func draw(ctx *core.Context) {
	screen := core.MustResource[*ebiten.Image](ctx)
	e := core.MustResource[*Extractor](ctx)

	// extraction
	for phase, fns := range e.extractors {
		buf := e.renderables[phase][:0]
		for _, fn := range fns {
			buf = append(buf, fn(ctx)...)
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
