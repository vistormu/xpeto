package shape

import (
	"image/color"

	"github.com/vistormu/xpeto/pkg/geometry"
	"github.com/vistormu/xpeto/pkg/render"
)

type Line struct {
	geometry.Line[float32]
	Shape
}

func NewLine(p0, p1 geometry.Vector[float32]) Line {
	return Line{
		Line:  geometry.NewLine(p0, p1),
		Shape: newShape(),
	}
}

func (l *Line) AddStroke(c color.Color, w float32) { l.Shape.AddStroke(c, w) }
func (l *Line) AddOrder(layer, order uint16)       { l.Shape.AddOrder(layer, order) }
func (l *Line) SetAnchor(a render.Anchor)          { l.Shape.Anchor = a }
