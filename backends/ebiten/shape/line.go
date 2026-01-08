package shape

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/go-dsa/geometry"
	"github.com/vistormu/xpeto/backends/ebiten/shared"
)

func drawLine(screen *ebiten.Image, r renderable) {
	s := r.Shape
	l := r.Shape.Line
	tr := r.Transform

	start := l.Point
	end := geometry.Vector[float32]{
		X: start.X + l.Direction.X,
		Y: start.Y + l.Direction.Y,
	}

	minX := min(start.X, end.X)
	minY := min(start.Y, end.Y)
	maxX := max(start.X, end.X)
	maxY := max(start.Y, end.Y)

	bw := float64(maxX - minX)
	bh := float64(maxY - minY)

	ax, ay := shared.Offset(bw, bh, r.anchor)

	tlx := tr.X + ax
	tly := tr.Y + ay

	x0 := tlx + float64(start.X-minX)
	y0 := tly + float64(start.Y-minY)
	x1 := tlx + float64(end.X-minX)
	y1 := tly + float64(end.Y-minY)

	if r.snap {
		x0 = math.Round(x0)
		y0 = math.Round(y0)
		x1 = math.Round(x1)
		y1 = math.Round(y1)
	}

	for _, s := range s.Strokes {
		if s.Width <= 0 {
			continue
		}
		vector.StrokeLine(
			screen,
			float32(x0), float32(y0),
			float32(x1), float32(y1),
			s.Width,
			s.Color,
			r.antialias,
		)
	}
}
