package shape

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/xpeto/backends/ebiten/shared"
	"github.com/vistormu/xpeto/pkg/shape"
)

func drawArrow(screen *ebiten.Image, r renderable) {
	s := r.Shape
	a := r.Shape.Arrow
	tr := r.Transform

	x0 := float64(a.Start.X)
	y0 := float64(a.Start.Y)
	x1 := float64(a.End.X)
	y1 := float64(a.End.Y)

	dx := x1 - x0
	dy := y1 - y0
	len2 := dx*dx + dy*dy
	if len2 == 0 {
		return
	}

	invLen := 1.0 / math.Sqrt(len2)
	ux := dx * invLen
	uy := dy * invLen

	hl := float64(a.HeadLength)
	hw := float64(a.HeadWidth)

	shaftLen := math.Sqrt(len2)
	if hl <= 0 {
		hl = math.Min(12, shaftLen*0.25)
	}
	if hl > shaftLen {
		hl = shaftLen
	}
	if hw <= 0 {
		hw = hl * 0.75
	}

	bx := x1 - ux*hl
	by := y1 - uy*hl

	px := -uy
	py := ux

	hx0 := bx + px*(hw*0.5)
	hy0 := by + py*(hw*0.5)
	hx1 := bx - px*(hw*0.5)
	hy1 := by - py*(hw*0.5)

	minX := min(min(x0, x1), min(hx0, hx1))
	minY := min(min(y0, y1), min(hy0, hy1))
	maxX := max(max(x0, x1), max(hx0, hx1))
	maxY := max(max(y0, y1), max(hy0, hy1))

	bw := maxX - minX
	bh := maxY - minY

	ax, ay := shared.Offset(bw, bh, r.anchor)

	tlx := tr.X + ax
	tly := tr.Y + ay

	wsx := tlx + (x0 - minX)
	wsy := tly + (y0 - minY)

	wtx := tlx + (x1 - minX)
	wty := tly + (y1 - minY)

	wbx0 := tlx + (hx0 - minX)
	wby0 := tly + (hy0 - minY)

	wbx1 := tlx + (hx1 - minX)
	wby1 := tly + (hy1 - minY)

	if r.snap {
		wsx = math.Round(wsx)
		wsy = math.Round(wsy)
		wtx = math.Round(wtx)
		wty = math.Round(wty)
		wbx0 = math.Round(wbx0)
		wby0 = math.Round(wby0)
		wbx1 = math.Round(wbx1)
		wby1 = math.Round(wby1)
	}

	var head vector.Path
	head.MoveTo(float32(wtx), float32(wty))   // tip
	head.LineTo(float32(wbx0), float32(wby0)) // base corner
	head.LineTo(float32(wbx1), float32(wby1)) // base corner
	head.Close()

	for _, f := range s.Fills {
		switch f.Type {
		case shape.FillSolid:
			var cs ebiten.ColorScale
			cs.ScaleWithColor(f.Color)

			draw := &vector.DrawPathOptions{
				AntiAlias:  r.antialias,
				ColorScale: cs,
			}
			vector.FillPath(screen, &head, nil, draw)
		default:
			// TODO: gradients / image fills
		}
	}

	for _, s := range s.Strokes {
		if s.Width <= 0 {
			continue
		}

		wbx := tlx + (bx - minX)
		wby := tly + (by - minY)
		if r.snap {
			wbx = math.Round(wbx)
			wby = math.Round(wby)
		}

		vector.StrokeLine(
			screen,
			float32(wsx), float32(wsy),
			float32(wbx), float32(wby),
			s.Width,
			s.Color,
			r.antialias,
		)

		var cs ebiten.ColorScale
		cs.ScaleWithColor(s.Color)

		draw := &vector.DrawPathOptions{
			AntiAlias:  r.antialias,
			ColorScale: cs,
		}

		st := &vector.StrokeOptions{Width: s.Width}
		vector.StrokePath(screen, &head, st, draw)
	}
}
