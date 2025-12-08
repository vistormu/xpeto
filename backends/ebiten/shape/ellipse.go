package shape

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/pkg/shape"
	"github.com/vistormu/xpeto/pkg/transform"
)

type ellipse struct {
	shape.Ellipse
	transform.Transform
}

func extractEllipse(w *ecs.World) []ellipse {
	q := ecs.NewQuery2[shape.Ellipse, transform.Transform](w)

	out := make([]ellipse, 0)

	for _, b := range q.Iter() {
		e, t := b.Components()

		out = append(out, ellipse{*e, *t})
	}

	return out
}

func sortEllipse(e ellipse) uint64 {
	return (uint64(e.Layer) << 16) | uint64(e.Order)
}

func drawEllipse(screen *ebiten.Image, e ellipse) {
	if e.RadiusX <= 0 || e.RadiusY <= 0 {
		return
	}

	// fill
	for _, f := range e.Fill {
		if !f.Visible {
			continue
		}
		vector.FillCircle(screen, float32(e.X), float32(e.Y), e.RadiusX, f.Color, false)
	}

	// stroke
	for _, s := range e.Stroke {
		if !s.Visible || s.Width <= 0 {
			continue
		}
		vector.StrokeCircle(screen, float32(e.X), float32(e.Y), e.RadiusX, s.Width, s.Color, false)
	}
}
