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

func (r *Rect) AddFillSolid(c color.Color) {
	r.Shape.AddFillSolid(c)
}

func (r *Rect) AddStroke(c color.Color, w float32) {
	r.Shape.AddStroke(c, w)
}

func (r *Rect) AddOrder(l, o uint16) {
	r.Shape.AddOrder(l, o)
}
