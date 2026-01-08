package shape

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/xpeto/backends/ebiten/shared"
	"github.com/vistormu/xpeto/pkg/shape"
)

func drawCapsule(screen *ebiten.Image, r renderable) {
	sh := r.Shape
	c := r.Shape.Capsule
	tr := r.Transform

	seg := c.Segment
	rad := float64(c.Radius)
	if rad <= 0 {
		return
	}

	start := seg.Start
	end := seg.End

	minX := min(start.X, end.X) - float32(rad)
	minY := min(start.Y, end.Y) - float32(rad)
	maxX := max(start.X, end.X) + float32(rad)
	maxY := max(start.Y, end.Y) + float32(rad)

	bw := float64(maxX - minX)
	bh := float64(maxY - minY)

	ax, ay := shared.Offset(bw, bh, r.anchor)

	tlx := tr.X + ax - float64(minX)
	tly := tr.Y + ay - float64(minY)

	x0 := tlx + float64(start.X)
	y0 := tly + float64(start.Y)
	x1 := tlx + float64(end.X)
	y1 := tly + float64(end.Y)

	if r.snap {
		x0 = math.Round(x0)
		y0 = math.Round(y0)
		x1 = math.Round(x1)
		y1 = math.Round(y1)
	}

	for _, st := range sh.Strokes {
		if st.Width <= 0 {
			continue
		}

		vector.StrokeLine(
			screen,
			float32(x0), float32(y0),
			float32(x1), float32(y1),
			st.Width+float32(rad*2),
			st.Color,
			r.antialias,
		)
	}

	for _, f := range sh.Fills {
		switch f.Type {
		case shape.FillSolid:
			vector.StrokeLine(
				screen,
				float32(x0), float32(y0),
				float32(x1), float32(y1),
				float32(rad*2),
				f.Color,
				r.antialias,
			)

			vector.FillCircle(screen, float32(x0), float32(y0), float32(rad), f.Color, r.antialias)
			vector.FillCircle(screen, float32(x1), float32(y1), float32(rad), f.Color, r.antialias)
		default:
			// not implemented in v0.1.0
		}
	}

}
