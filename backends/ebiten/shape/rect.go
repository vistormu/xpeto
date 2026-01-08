package shape

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/xpeto/backends/ebiten/shared"
	"github.com/vistormu/xpeto/pkg/shape"
)

func drawRect(screen *ebiten.Image, re renderable) {
	s := re.Shape
	r := re.Shape.Rect
	tr := re.Transform

	if r.Width <= 0 || r.Height <= 0 {
		return
	}

	bw := float64(r.Width)
	bh := float64(r.Height)
	ax, ay := shared.Offset(bw, bh, re.anchor)

	x := tr.X + ax
	y := tr.Y + ay

	if re.snap {
		x = math.Round(x)
		y = math.Round(y)
	}

	xf := float32(x)
	yf := float32(y)

	for _, f := range s.Fills {
		switch f.Type {
		case shape.FillSolid:
			vector.FillRect(screen, xf, yf, r.Width, r.Height, f.Color, re.antialias)
		default:
			// TODO: gradients / image fills
		}
	}

	for _, s := range s.Strokes {
		if s.Width <= 0 {
			continue
		}

		sx, sy := xf, yf
		if re.snap && int(s.Width)%2 == 1 {
			sx += 0.5
			sy += 0.5
		}

		vector.StrokeRect(screen, sx, sy, r.Width, r.Height, s.Width, s.Color, re.antialias)
	}
}
