package shape

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/pkg/geometry"
	"github.com/vistormu/xpeto/pkg/transform"
)

// =========
// component
// =========
type Path struct {
	geometry.Path[float32]
	Shape
}

func NewPath() Path {
	return Path{
		Path:  geometry.NewPath[float32](),
		Shape: newShape(),
	}
}

func (p Path) AddStroke(c color.Color, w float32) Path {
	p.Shape = p.Shape.AddStroke(c, w)
	return p
}

func (p Path) AddOrder(l, o uint16) Path {
	p.Shape = p.Shape.AddOrder(l, o)
	return p
}

// ==========
// renderable
// ==========
type path struct {
	Path
	transform.Transform
}

func extractPath(w *ecs.World) []path {
	q := ecs.NewQuery2[Path, transform.Transform](w)

	out := make([]path, 0)

	for _, b := range q.Iter() {
		p, t := b.Components()

		out = append(out, path{*p, *t})
	}

	return out
}

func sortPath(p path) uint64 {
	return (uint64(p.Layer) << 16) | uint64(p.Order)
}

func drawPath(screen *ebiten.Image, p path) {
	if len(p.Points) < 2 {
		return
	}

	var vPath vector.Path

	first := p.Points[0]
	vPath.MoveTo(first.X, first.Y)

	for _, point := range p.Points[1:] {
		vPath.LineTo(point.X, point.Y)
	}

	for _, s := range p.Stroke {
		if !s.Visible || s.Width <= 0 {
			continue
		}

		opts := &vector.StrokeOptions{
			Width: s.Width,
		}

		cs := ebiten.ColorScale{}
		cs.ScaleWithColor(s.Color)

		draw := &vector.DrawPathOptions{
			AntiAlias:  true,
			ColorScale: cs,
		}

		vector.StrokePath(screen, &vPath, opts, draw)
	}
}
