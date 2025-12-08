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

func (e Ellipse) AddFillSolid(c color.Color) Ellipse {
	e.Shape = e.Shape.AddFillSolid(c)
	return e
}

func (e Ellipse) AddStroke(c color.Color, w float32) Ellipse {
	e.Shape = e.Shape.AddStroke(c, w)
	return e
}

func (e Ellipse) AddOrder(l, o uint16) Ellipse {
	e.Shape = e.Shape.AddOrder(l, o)
	return e
}
