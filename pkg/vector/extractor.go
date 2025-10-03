package vector

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"

	"github.com/vistormu/xpeto/pkg/render"
	"github.com/vistormu/xpeto/pkg/transform"
)

func packKey(layer, order uint16) uint64 {
	return (uint64(layer) << 16) | uint64(order)
}

func extractCircles(ctx *core.Context) []render.Renderable {
	w := core.MustResource[*ecs.World](ctx)
	entities := w.Query(ecs.And(
		ecs.Has[*Circle](),
		ecs.Has[*transform.Transform](),
	))

	renderables := make([]render.Renderable, 0)
	for _, e := range entities {
		// components
		transform, _ := ecs.GetComponent[*transform.Transform](w, e)
		circle, _ := ecs.GetComponent[*Circle](w, e)

		item := &circleItem{
			radius:    circle.Radius,
			fill:      circle.Fill,
			stroke:    circle.Stroke,
			linewidth: circle.Linewidth,
			x:         transform.Position.X,
			y:         transform.Position.Y,
			layer:     circle.Layer,
			order:     circle.Order,
			key:       packKey(circle.Layer, circle.Order),
		}
		renderables = append(renderables, item)
	}

	return renderables
}

func extractRects(ctx *core.Context) []render.Renderable {
	w := core.MustResource[*ecs.World](ctx)
	entities := w.Query(ecs.And(
		ecs.Has[*Rect](),
		ecs.Has[*transform.Transform](),
	))

	renderables := make([]render.Renderable, 0)
	for _, e := range entities {
		// components
		transform, _ := ecs.GetComponent[*transform.Transform](w, e)
		rect, _ := ecs.GetComponent[*Rect](w, e)

		item := &rectItem{
			width:     rect.Width,
			height:    rect.Height,
			fill:      rect.Fill,
			stroke:    rect.Stroke,
			linewidth: rect.Linewidth,
			x:         transform.Position.X,
			y:         transform.Position.Y,
			layer:     rect.Layer,
			order:     rect.Order,
			key:       packKey(rect.Layer, rect.Order),
		}
		renderables = append(renderables, item)
	}

	return renderables
}
