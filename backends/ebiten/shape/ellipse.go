package shape

import (
	// "image"
	// "image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/xpeto/backends/ebiten/shared"
	"github.com/vistormu/xpeto/pkg/shape"
)

func drawEllipse(screen *ebiten.Image, r renderable) {
	s := r.Shape
	e := r.Shape.Ellipse
	tr := r.Transform

	rx := float64(e.RadiusX)
	ry := float64(e.RadiusY)
	if rx <= 0 || ry <= 0 {
		return
	}

	w := rx * 2
	h := ry * 2
	dx, dy := shared.Offset(w, h, r.anchor)

	cx := tr.X + dx + rx
	cy := tr.Y + dy + ry

	if r.snap {
		cx = math.Round(cx)
		cy = math.Round(cy)
	}

	const segments = 64

	var p vector.Path
	p.MoveTo(float32(cx+rx), float32(cy))
	for i := 1; i <= segments; i++ {
		theta := (2 * math.Pi * float64(i)) / float64(segments)
		x := cx + math.Cos(theta)*rx
		y := cy + math.Sin(theta)*ry
		p.LineTo(float32(x), float32(y))
	}
	p.Close()

	// fill
	for _, f := range s.Fills {
		switch f.Type {
		case shape.FillSolid:
			var dop vector.DrawPathOptions
			dop.AntiAlias = r.antialias
			dop.ColorScale.ScaleWithColor(f.Color)

			vector.FillPath(screen, &p, nil, &dop)
		default:
			// TODO: gradients / image fills
		}
	}

	// stroke
	for _, s := range s.Strokes {
		if s.Width <= 0 {
			continue
		}

		var sop vector.StrokeOptions
		sop.Width = s.Width

		var dop vector.DrawPathOptions
		dop.AntiAlias = r.antialias
		dop.ColorScale.ScaleWithColor(s.Color)

		vector.StrokePath(screen, &p, &sop, &dop)
	}
}

// func emptySubImage(screen *ebiten.Image) *ebiten.Image {
// 	return screen.SubImage(imageRect1x1()).(*ebiten.Image)
// }

// func imageRect1x1() (r image.Rectangle) {
// 	return image.Rect(0, 0, 1, 1)
// }

// func colorToFloats(c color.Color) (r, g, b, a float32) {
// 	rr, gg, bb, aa := c.RGBA()
// 	const inv = 1.0 / 65535.0
// 	return float32(float64(rr) * inv),
// 		float32(float64(gg) * inv),
// 		float32(float64(bb) * inv),
// 		float32(float64(aa) * inv)
// }
