package shape

import (
	"image/color"

	"github.com/vistormu/xpeto/pkg/geometry"
	"github.com/vistormu/xpeto/pkg/render"
)

type Segment struct {
	geometry.Segment[float32]
	Shape
}

func NewSegment(start, end geometry.Vector[float32]) Segment {
	return Segment{
		Segment: geometry.NewSegment(start, end),
		Shape:   newShape(),
	}
}

func (s *Segment) AddStroke(c color.Color, w float32) { s.Shape.AddStroke(c, w) }
func (s *Segment) AddOrder(layer, order uint16)       { s.Shape.AddOrder(layer, order) }
func (s *Segment) SetAnchor(a render.Anchor)          { s.Shape.Anchor = a }
