package shape

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/pkg/shape"
	"github.com/vistormu/xpeto/pkg/transform"
)

type rect struct {
	shape.Rect
	transform.Transform
}

func extractRect(w *ecs.World) []rect {
	q := ecs.NewQuery2[shape.Rect, transform.Transform](w)

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
