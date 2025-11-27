package geometry

import (
	c "github.com/vistormu/go-dsa/constraints"
)

type Ellipse[T c.Number] struct {
	RadiusX T
	RadiusY T
}

func NewEllipse[T c.Number](rx, ry T) Ellipse[T] {
	return Ellipse[T]{rx, ry}
}

func NewEllipseHW[T c.Number](h, w T) Ellipse[T] {
	return Ellipse[T]{w / 2, h / 2}
}

func NewCircle[T c.Number](r T) Ellipse[T] {
	return Ellipse[T]{r, r}
}

func NewCircleD[T c.Number](d T) Ellipse[T] {
	return Ellipse[T]{d / 2, d / 2}
}
