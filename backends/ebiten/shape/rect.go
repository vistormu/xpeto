package shape

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/xpeto/backends/ebiten/shared"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/window"
	"github.com/vistormu/xpeto/pkg/shape"
	"github.com/vistormu/xpeto/pkg/transform"
)

type rect struct {
	shape.Rect
	transform.Transform

	snap      bool
	antialias bool
}

func extractRect(w *ecs.World) []rect {
	q := ecs.NewQuery2[shape.Rect, transform.Transform](w)

	sc, _ := ecs.GetResource[window.Scaling](w)
	rw, _ := ecs.GetResource[window.RealWindow](w)

	out := make([]rect, 0)
	for _, b := range q.Iter() {
		r, t := b.Components()
		if !r.Visible {
			continue
		}
		out = append(out, rect{
			Rect:      *r,
			Transform: *t,
			snap:      sc.SnapPixels,
			antialias: rw.AntiAliasing,
		})
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

	bw := float64(r.Width)
	bh := float64(r.Height)
	ax, ay := shared.Offset(bw, bh, r.Anchor)

	x := r.X + ax
	y := r.Y + ay

	if r.snap {
		x = math.Round(x)
		y = math.Round(y)
	}

	xf := float32(x)
	yf := float32(y)

	for _, f := range r.Fill {
		if !f.Visible {
			continue
		}
		switch f.Type {
		case shape.FillSolid:
			vector.FillRect(screen, xf, yf, r.Width, r.Height, f.Color, r.antialias)
		default:
			// TODO: gradients / image fills
		}
	}

	for _, s := range r.Stroke {
		if !s.Visible || s.Width <= 0 {
			continue
		}

		sx, sy := xf, yf
		if r.snap && int(s.Width)%2 == 1 {
			sx += 0.5
			sy += 0.5
		}

		vector.StrokeRect(screen, sx, sy, r.Width, r.Height, s.Width, s.Color, r.antialias)
	}
}
