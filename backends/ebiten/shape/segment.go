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

type segment struct {
	shape.Segment
	transform.Transform

	snap      bool
	antialias bool
}

func extractSegment(w *ecs.World) []segment {
	q := ecs.NewQuery2[shape.Segment, transform.Transform](w)

	sc, _ := ecs.GetResource[window.Scaling](w)
	rw, _ := ecs.GetResource[window.RealWindow](w)

	out := make([]segment, 0)
	for _, b := range q.Iter() {
		s, t := b.Components()
		if !s.Visible {
			continue
		}
		out = append(out, segment{
			Segment:   *s,
			Transform: *t,
			snap:      sc.SnapPixels,
			antialias: rw.AntiAliasing,
		})
	}
	return out
}

func sortSegment(s segment) uint64 {
	return (uint64(s.Layer) << 16) | uint64(s.Order)
}

func drawSegment(screen *ebiten.Image, s segment) {
	start := s.Start
	end := s.End

	minX := min(start.X, end.X)
	minY := min(start.Y, end.Y)
	maxX := max(start.X, end.X)
	maxY := max(start.Y, end.Y)

	bw := float64(maxX - minX)
	bh := float64(maxY - minY)

	ax, ay := shared.Offset(bw, bh, s.Anchor)

	tlx := s.X + ax - float64(minX)
	tly := s.Y + ay - float64(minY)

	x0 := tlx + float64(start.X)
	y0 := tly + float64(start.Y)
	x1 := tlx + float64(end.X)
	y1 := tly + float64(end.Y)

	if s.snap {
		x0 = math.Round(x0)
		y0 = math.Round(y0)
		x1 = math.Round(x1)
		y1 = math.Round(y1)
	}

	for _, st := range s.Stroke {
		if !st.Visible || st.Width <= 0 {
			continue
		}
		vector.StrokeLine(
			screen,
			float32(x0), float32(y0),
			float32(x1), float32(y1),
			st.Width,
			st.Color,
			s.antialias,
		)
	}
}
