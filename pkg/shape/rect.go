package shape

import (
	"image/color"

	"github.com/vistormu/xpeto/pkg/geometry"
)

type Rect struct {
	geometry.Rect[float32]
	Shape
}

func NewRect(w, h float32) Rect {
	return Rect{
		Rect:  geometry.NewRect(w, h),
		Shape: newShape(),
	}
}

func (r Rect) AddFillSolid(c color.Color) Rect {
	r.Shape = r.Shape.AddFillSolid(c)
	return r
}

func (r Rect) AddStroke(c color.Color, w float32) Rect {
	r.Shape = r.Shape.AddStroke(c, w)
	return r
}

func (r Rect) AddOrder(l, o uint16) Rect {
	r.Shape = r.Shape.AddOrder(l, o)
	return r
}
