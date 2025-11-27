package geometry

import (
	c "github.com/vistormu/go-dsa/constraints"
)

type Rect[T c.Number] struct {
	Width  T
	Height T
}

func NewRect[T c.Number](w, h T) Rect[T] {
	return Rect[T]{w, h}
}

func NewSquare[T c.Number](r T) Rect[T] {
	return Rect[T]{r, r}
}
