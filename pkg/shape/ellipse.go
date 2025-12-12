package shape

import (
	"image/color"

	"github.com/vistormu/xpeto/pkg/geometry"
)

type Ellipse struct {
	geometry.Ellipse[float32]
	Shape
}

func NewCircle(r float32) Ellipse {
	return Ellipse{
		Ellipse: geometry.NewCircle(r),
		Shape:   newShape(),
	}
}

func (e *Ellipse) AddFillSolid(c color.Color) {
	e.Shape.AddFillSolid(c)
}

func (e *Ellipse) AddStroke(c color.Color, w float32) {
	e.Shape.AddStroke(c, w)
}

func (e *Ellipse) AddOrder(l, o uint16) {
	e.Shape.AddOrder(l, o)
}
