package geometry

import (
	c "github.com/vistormu/go-dsa/constraints"
)

type Ellipse[T c.Number] struct {
	RadiusX T
	RadiusY T
}

func (s Ellipse[T]) Bounds() Rect[T] {
	return Rect[T]{}
}

func NewEllipse[T c.Number](rx, ry T) Geometry[T] {
	return Geometry[T]{
		Kind:    GeometryEllipse,
		Ellipse: Ellipse[T]{rx, ry},
	}
}

func NewEllipseHW[T c.Number](h, w T) Geometry[T] {
	return Geometry[T]{
		Kind:    GeometryEllipse,
		Ellipse: Ellipse[T]{w / 2, h / 2},
	}
}

func NewCircle[T c.Number](r T) Geometry[T] {
	return Geometry[T]{
		Kind:    GeometryEllipse,
		Ellipse: Ellipse[T]{r, r},
	}
}

func NewCircleD[T c.Number](d T) Geometry[T] {
	return Geometry[T]{
		Kind:    GeometryEllipse,
		Ellipse: Ellipse[T]{d / 2, d / 2},
	}
}
