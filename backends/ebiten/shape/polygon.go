package shape

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/xpeto/backends/ebiten/shared"
	"github.com/vistormu/xpeto/pkg/shape"
)

func drawPolygon(screen *ebiten.Image, r renderable) {
	sh := r.Shape
	pg := r.Shape.Polygon
	tr := r.Transform

	pts := pg.Points
	if len(pts) < 3 {
		return
	}

	minX, maxX := pts[0].X, pts[0].X
	minY, maxY := pts[0].Y, pts[0].Y
	for _, pt := range pts[1:] {
		minX = min(minX, pt.X)
		minY = min(minY, pt.Y)
		maxX = max(maxX, pt.X)
		maxY = max(maxY, pt.Y)
	}

	bw := float64(maxX - minX)
	bh := float64(maxY - minY)

	ax, ay := shared.Offset(bw, bh, r.anchor)
	tlx := tr.X + ax - float64(minX)
	tly := tr.Y + ay - float64(minY)

	var path vector.Path
	path.MoveTo(float32(tlx+float64(pts[0].X)), float32(tly+float64(pts[0].Y)))
	for _, pt := range pts[1:] {
		x := tlx + float64(pt.X)
		y := tly + float64(pt.Y)
		if r.snap {
			x = math.Round(x)
			y = math.Round(y)
		}
		path.LineTo(float32(x), float32(y))
	}
	path.Close()

	for _, f := range sh.Fills {
		switch f.Type {
		case shape.FillSolid:
			var dop vector.DrawPathOptions
			dop.AntiAlias = r.antialias
			dop.ColorScale.ScaleWithColor(f.Color)
			vector.FillPath(screen, &path, nil, &dop)
		default:
			// TODO: gradients / image fills
		}
	}

	for _, st := range sh.Strokes {
		if st.Width <= 0 {
			continue
		}

		var sop vector.StrokeOptions
		sop.Width = st.Width

		var dop vector.DrawPathOptions
		dop.AntiAlias = r.antialias
		dop.ColorScale.ScaleWithColor(st.Color)

		vector.StrokePath(screen, &path, &sop, &dop)
	}
}
