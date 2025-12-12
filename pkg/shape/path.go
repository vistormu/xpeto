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

func (p *Path) AddStroke(c color.Color, w float32) {
	p.Shape.AddStroke(c, w)
}

// func (p *Path) SetColor(c color.Color)

func (p *Path) AddOrder(l, o uint16) {
	p.Shape.AddOrder(l, o)
}
