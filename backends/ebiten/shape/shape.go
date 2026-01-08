package shape

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/window"
	"github.com/vistormu/xpeto/pkg/render"
	"github.com/vistormu/xpeto/pkg/shape"
	"github.com/vistormu/xpeto/pkg/transform"
)

type renderable struct {
	shape.Shape
	transform.Transform
	anchor    render.AnchorType
	snap      bool
	antialias bool
}

type renderableBuffer struct {
	renderables []renderable
}

func newRenderablebuffer() renderableBuffer {
	return renderableBuffer{
		renderables: make([]renderable, 0),
	}
}

func extractShape(w *ecs.World) []renderable {
	buf := ecs.EnsureResource(w, newRenderablebuffer)
	buf.renderables = buf.renderables[:0]

	sc, ok := ecs.GetResource[window.Scaling](w)
	if !ok {
		return buf.renderables
	}
	rw, ok := ecs.GetResource[window.RealWindow](w)
	if !ok {
		return buf.renderables
	}

	q := ecs.NewQuery2[shape.Shape, transform.Transform](w)

	for _, b := range q.Iter() {
		s, t := b.Components()
		e := b.Entity()

		anchor := render.AnchorCenter
		an, ok := ecs.GetComponent[render.Anchor](w, e)
		if ok {
			anchor = an.Type
		}

		buf.renderables = append(buf.renderables, renderable{
			Shape:     *s,
			Transform: *t,
			anchor:    anchor,
			snap:      sc.SnapPixels,
			antialias: rw.AntiAliasing,
		})
	}

	return buf.renderables
}

func sortShape(r renderable) uint64 {
	return uint64(r.OrderKey)
}

func drawShape(screen *ebiten.Image, r renderable) {
	switch r.Kind {
	case shape.Arrow:
		drawArrow(screen, r)
	case shape.Capsule:
		drawCapsule(screen, r)
	case shape.Ellipse:
		drawEllipse(screen, r)
	case shape.Line:
		drawLine(screen, r)
	case shape.Path:
		drawPath(screen, r)
	case shape.Polygon:
		drawPolygon(screen, r)
	case shape.Ray:
		drawRay(screen, r)
	case shape.Rect:
		drawRect(screen, r)
	case shape.Segment:
		drawSegment(screen, r)
	}
}
