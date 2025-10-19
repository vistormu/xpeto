package vector

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/pkg/transform"
	"github.com/vistormu/xpeto/pkg/render"
)

// ====
// item
// ====
type Circle struct {
	Radius    float32
	Fill      color.Color
	Stroke    color.Color
	Linewidth float32
	Layer     uint16
	Order     uint16
}

type circleItem struct {
	radius    float32
	fill      color.Color
	stroke    color.Color
	linewidth float32
	x, y      float32
	layer     uint16
	order     uint16
	key       uint64
}

func (ci *circleItem) Draw(screen *ebiten.Image) {
	if ci.fill != nil && ci.radius > 0 {
		vector.FillCircle(screen, ci.x, ci.y, ci.radius, ci.fill, false)
	}
	if ci.linewidth > 0 && ci.stroke != nil && ci.radius > 0 {
		vector.StrokeCircle(screen, ci.x, ci.y, ci.radius, ci.linewidth, ci.stroke, false)
	}
}

func (ci *circleItem) SortKey() uint64 {
	return ci.key
}

// =========
// extractor
// =========
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
