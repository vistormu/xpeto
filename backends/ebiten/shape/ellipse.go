package shape

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/xpeto/backends/ebiten/shared"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/window"
	"github.com/vistormu/xpeto/pkg/shape"
	"github.com/vistormu/xpeto/pkg/transform"
)

type ellipse struct {
	shape.Ellipse
	transform.Transform
	snap      bool
	antialias bool
}

func extractEllipse(w *ecs.World) []ellipse {
	q := ecs.NewQuery2[shape.Ellipse, transform.Transform](w)

	sc, _ := ecs.GetResource[window.Scaling](w)
	rw, _ := ecs.GetResource[window.RealWindow](w)

	out := make([]ellipse, 0)
	for _, b := range q.Iter() {
		e, t := b.Components()

		if !e.Visible {
			continue
		}

		out = append(out, ellipse{
			Ellipse:   *e,
			Transform: *t,
			snap:      sc.SnapPixels,
			antialias: rw.AntiAliasing,
		})
	}
	return out
}

func sortEllipse(e ellipse) uint64 {
	return (uint64(e.Layer) << 16) | uint64(e.Order)
}

func drawEllipse(screen *ebiten.Image, e ellipse) {
	rx := float64(e.RadiusX)
	ry := float64(e.RadiusY)
	if rx <= 0 || ry <= 0 {
		return
	}

	w := rx * 2
	h := ry * 2
	dx, dy := shared.Offset(w, h, e.Anchor)

	cx := e.X + dx + rx
	cy := e.Y + dy + ry

	if e.snap {
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
	for _, f := range e.Fill {
		if !f.Visible {
			continue
		}
		switch f.Type {
		case shape.FillSolid:
			var dop vector.DrawPathOptions
			dop.AntiAlias = e.antialias
			dop.ColorScale.ScaleWithColor(f.Color)

			vector.FillPath(screen, &p, nil, &dop)
		default:
			// TODO: gradients / image fills
		}
	}

	// stroke
	for _, s := range e.Stroke {
		if !s.Visible || s.Width <= 0 {
			continue
		}

		var sop vector.StrokeOptions
		sop.Width = s.Width

		var dop vector.DrawPathOptions
		dop.AntiAlias = e.antialias
		dop.ColorScale.ScaleWithColor(s.Color)

		vector.StrokePath(screen, &p, &sop, &dop)
	}
}

func emptySubImage(screen *ebiten.Image) *ebiten.Image {
	return screen.SubImage(imageRect1x1()).(*ebiten.Image)
}

func imageRect1x1() (r image.Rectangle) {
	return image.Rect(0, 0, 1, 1)
}

func colorToFloats(c color.Color) (r, g, b, a float32) {
	rr, gg, bb, aa := c.RGBA()
	const inv = 1.0 / 65535.0
	return float32(float64(rr) * inv),
		float32(float64(gg) * inv),
		float32(float64(bb) * inv),
		float32(float64(aa) * inv)
}
