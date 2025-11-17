package shape

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/xpeto/core/ecs"
	g "github.com/vistormu/xpeto/pkg/geometry"
	"github.com/vistormu/xpeto/pkg/transform"
)

type ellipse struct {
	radiusX   float32
	radiusY   float32
	fill      color.Color
	stroke    color.Color
	linewidth float32
	x, y      float32
	layer     uint16
	order     uint16
}

func extractEllipse(w *ecs.World) []ellipse {
	q := ecs.NewQuery2[Shape, transform.Transform](w)

	out := make([]ellipse, 0)

	for _, b := range q.Iter() {
		s, t := b.Components()

		if s.Shape.Kind != g.GeometryEllipse {
			continue
		}

		out = append(out, ellipse{
			radiusX:   s.Shape.Ellipse.RadiusX,
			radiusY:   s.Shape.Ellipse.RadiusY,
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

func sortEllipse(e ellipse) uint64 {
	return (uint64(e.layer) << 16) | uint64(e.order)
}

func drawEllipse(screen *ebiten.Image, e ellipse) {
	if e.radiusX <= 0 || e.radiusY <= 0 {
		return
	}

	// fill
	if e.fill != nil {
		vector.FillCircle(screen, e.x, e.y, e.radiusX, e.fill, false)
	}

	// stroke
	if e.linewidth >= 0 && e.stroke != nil {
		vector.StrokeCircle(screen, e.x, e.y, e.radiusX, e.linewidth, e.stroke, false)
	}
}
