package geometry

import (
	c "github.com/vistormu/go-dsa/constraints"
)

type Rect[T c.Number] struct {
	Width  T
	Height T
}

func (s Rect[T]) Bounds() Rect[T] {
	return s
}

func NewRect[T c.Number](w, h T) Geometry[T] {
	return Geometry[T]{
		Kind: GeometryRect,
		Rect: Rect[T]{w, h},
	}
}

func NewSquare[T c.Number](r T) Geometry[T] {
	return Geometry[T]{
		Kind: GeometryRect,
		Rect: Rect[T]{r, r},
	}
}
