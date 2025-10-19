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
type Rect struct {
	Width     float32
	Height    float32
	Fill      color.Color
	Stroke    color.Color
	Linewidth float32
	Layer     uint16
	Order     uint16
}

type rectItem struct {
	width     float32
	height    float32
	fill      color.Color
	stroke    color.Color
	linewidth float32
	x, y      float32
	layer     uint16
	order     uint16
	key       uint64
}

//	func (ri *rectItem) Draw(screen *ebiten.Image) {
//		if ri.fill != nil && ri.width > 0 && ri.height > 0 {
//			vector.FillRect(screen, ri.x, ri.y, ri.width, ri.height, ri.fill, false)
//		}
//		if ri.linewidth > 0 && ri.stroke != nil && ri.width > 0 && ri.height > 0 {
//			vector.StrokeRect(screen, ri.x, ri.y, ri.width, ri.height, ri.linewidth, ri.stroke, false)
//		}
//	}
func (ri *rectItem) Draw(screen *ebiten.Image) {
	if ri.width <= 0 || ri.height <= 0 {
		return
	}

	x := ri.x - ri.width*0.5
	y := ri.y - ri.height*0.5

	if ri.fill != nil {
		vector.FillRect(screen, x, y, ri.width, ri.height, ri.fill, false)
	}
	if ri.stroke != nil && ri.linewidth > 0 {
		// optional crisp offset for odd stroke widths
		if int(ri.linewidth)%2 == 1 {
			x += 0.5
			y += 0.5
		}
		vector.StrokeRect(screen, x, y, ri.width, ri.height, ri.linewidth, ri.stroke, false)
	}
}

func (ri *rectItem) SortKey() uint64 {
	return ri.key
}

// ========
// extactor
// ========
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
