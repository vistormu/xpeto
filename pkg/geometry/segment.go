package geometry

import (
	c "github.com/vistormu/go-dsa/constraints"
)

type Segment[T c.Number] struct {
	Start Vector[T]
	End   Vector[T]
}

func NewSegment[T c.Number](start, end Vector[T]) Segment[T] {
	return Segment[T]{Start: start, End: end}
}
