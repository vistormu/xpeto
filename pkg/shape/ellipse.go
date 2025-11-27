package shape

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/pkg/geometry"
	"github.com/vistormu/xpeto/pkg/transform"
)

// =========
// component
// =========
type Ellipse struct {
	geometry.Ellipse[float32]
	Shape
}

func NewCircle(r float32) Ellipse {
	return Ellipse{
		Ellipse: geometry.NewCircle(r),
		Shape:   newShape(),
	}
}

func (e Ellipse) AddFillSolid(c color.Color) Ellipse {
	e.Shape = e.Shape.AddFillSolid(c)
	return e
}

func (e Ellipse) AddStroke(c color.Color, w float32) Ellipse {
	e.Shape = e.Shape.AddStroke(c, w)
	return e
}

func (e Ellipse) AddOrder(l, o uint16) Ellipse {
	e.Shape = e.Shape.AddOrder(l, o)
	return e
}

// ==========
// renderable
// ==========
type ellipse struct {
	Ellipse
	transform.Transform
}

func extractEllipse(w *ecs.World) []ellipse {
	q := ecs.NewQuery2[Ellipse, transform.Transform](w)

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
