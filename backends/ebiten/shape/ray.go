package shape

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/go-dsa/geometry"
	"github.com/vistormu/xpeto/backends/ebiten/shared"
)

func drawRay(screen *ebiten.Image, r renderable) {
	sh := r.Shape
	ry := r.Shape.Ray
	tr := r.Transform

	start := ry.Origin
	end := geometry.Vector[float32]{
		X: start.X + ry.Direction.X,
		Y: start.Y + ry.Direction.Y,
	}

	minX := min(start.X, end.X)
	minY := min(start.Y, end.Y)
	maxX := max(start.X, end.X)
	maxY := max(start.Y, end.Y)

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
			st.Width,
			st.Color,
			r.antialias,
		)
	}
}
