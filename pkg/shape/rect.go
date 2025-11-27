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
type Rect struct {
	geometry.Rect[float32]
	Shape
}

func NewRect(w, h float32) Rect {
	return Rect{
		Rect:  geometry.NewRect(w, h),
		Shape: newShape(),
	}
}

func (r Rect) AddFillSolid(c color.Color) Rect {
	r.Shape = r.Shape.AddFillSolid(c)
	return r
}

func (r Rect) AddStroke(c color.Color, w float32) Rect {
	r.Shape = r.Shape.AddStroke(c, w)
	return r
}

func (r Rect) AddOrder(l, o uint16) Rect {
	r.Shape = r.Shape.AddOrder(l, o)
	return r
}

// ==========
// renderable
// ==========
type rect struct {
	Rect
	transform.Transform
}

func extractRect(w *ecs.World) []rect {
	q := ecs.NewQuery2[Rect, transform.Transform](w)

	out := make([]rect, 0)

	for _, b := range q.Iter() {
		r, t := b.Components()

		out = append(out, rect{*r, *t})
	}

	return out
}

func sortRect(r rect) uint64 {
	return (uint64(r.Layer) << 16) | uint64(r.Order)
}

func drawRect(screen *ebiten.Image, r rect) {
	if r.Width <= 0 || r.Height <= 0 {
		return
	}

	x := float32(r.X) - r.Width*0.5
	y := float32(r.Y) - r.Height*0.5

	// fill
	for _, f := range r.Fill {
		if !f.Visible {
			continue
		}
		vector.FillRect(screen, x, y, r.Width, r.Height, f.Color, false)
	}

	// stroke
	for _, s := range r.Stroke {
		if !s.Visible || s.Width <= 0 {
			continue
		}
		if int(s.Width)%2 == 1 {
			x += 0.5
			y += 0.5
		}
		vector.StrokeRect(screen, x, y, r.Width, r.Height, s.Width, s.Color, false)
	}
}
