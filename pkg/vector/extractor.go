package vector

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/pkg/transform"

	"github.com/vistormu/xpeto/pkg/render"
)

func packKey(layer, order uint16) uint64 {
	return (uint64(layer) << 16) | uint64(order)
}

func extractCircles(w *ecs.World) []render.Renderable {
	q := ecs.NewQuery2[Circle, transform.Transform](w)

	renderables := make([]render.Renderable, 0)
	for _, b := range q.Iter() {
		// components
		ci := b.A()
		tr := b.B()

		item := &circleItem{
			radius:    ci.Radius,
			fill:      ci.Fill,
			stroke:    ci.Stroke,
			linewidth: ci.Linewidth,
			x:         float32(tr.X),
			y:         float32(tr.Y),
			layer:     ci.Layer,
			order:     ci.Order,
			key:       packKey(ci.Layer, ci.Order),
		}
		renderables = append(renderables, item)
	}

	return renderables
}

func extractRects(w *ecs.World) []render.Renderable {
	q := ecs.NewQuery2[Rect, transform.Transform](w)

	renderables := make([]render.Renderable, 0)
	for _, b := range q.Iter() {
		// components
		rect := b.A()
		tr := b.B()

		item := &rectItem{
			width:     rect.Width,
			height:    rect.Height,
			fill:      rect.Fill,
			stroke:    rect.Stroke,
			linewidth: rect.Linewidth,
			x:         float32(tr.X),
			y:         float32(tr.Y),
			layer:     rect.Layer,
			order:     rect.Order,
			key:       packKey(rect.Layer, rect.Order),
		}
		renderables = append(renderables, item)
	}

	return renderables
}
