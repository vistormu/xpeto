package shape

import (
	"image/color"

	"github.com/vistormu/xpeto/pkg/geometry"
	"github.com/vistormu/xpeto/pkg/render"
)

type Arrow struct {
	geometry.Arrow[float32]
	Shape
}

func NewArrow(start, end geometry.Vector[float32], headLength, headWidth float32) Arrow {
	return Arrow{
		Arrow: geometry.NewArrow(start, end, headLength, headWidth),
		Shape: newShape(),
	}
}

func (a *Arrow) AddFillSolid(c color.Color)         { a.Shape.AddFillSolid(c) }
func (a *Arrow) AddStroke(c color.Color, w float32) { a.Shape.AddStroke(c, w) }
func (a *Arrow) AddOrder(layer, order uint16)       { a.Shape.AddOrder(layer, order) }
func (a *Arrow) SetAnchor(anchor render.Anchor)     { a.Shape.Anchor = anchor }
