package shape

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/pkg/geometry"
	"github.com/vistormu/xpeto/pkg/transform"
)

type rect struct {
	width     float32
	height    float32
	fill      color.Color
	stroke    color.Color
	linewidth float32
	x, y      float32
	layer     uint16
	order     uint16
}

func extractRect(w *ecs.World) []rect {
	q := ecs.NewQuery2[Shape, transform.Transform](w)

	out := make([]rect, 0)

	for _, b := range q.Iter() {
		s, t := b.Components()

		if s.Shape.Kind != geometry.GeometryRect {
			continue
		}

		out = append(out, rect{
			width:     s.Shape.Rect.Width,
			height:    s.Shape.Rect.Height,
			fill:      s.Fill.Color,
			stroke:    s.Stroke.Fill.Color,
			linewidth: s.Stroke.Width,
			x:         float32(t.X),
			y:         float32(t.Y),
			layer:     s.Layer,
			order:     s.Order,
		})
	}

	return out
}

func sortRect(r rect) uint64 {
	return (uint64(r.layer) << 16) | uint64(r.order)
}

func drawRect(screen *ebiten.Image, r rect) {
	if r.width <= 0 || r.height <= 0 {
		return
	}

	x := r.x - r.width*0.5
	y := r.y - r.height*0.5

	if r.fill != nil {
		vector.FillRect(screen, x, y, r.width, r.height, r.fill, false)
	}
	if r.stroke != nil && r.linewidth > 0 {
		// optional crisp offset for odd stroke widths
		if int(r.linewidth)%2 == 1 {
			x += 0.5
			y += 0.5
		}
		vector.StrokeRect(screen, x, y, r.width, r.height, r.linewidth, r.stroke, false)
	}
}
