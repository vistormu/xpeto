package shape

import (
	"image/color"

	"github.com/vistormu/xpeto/pkg/geometry"
)

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
