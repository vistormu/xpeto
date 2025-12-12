package geometry

import (
	c "github.com/vistormu/go-dsa/constraints"
)

type Arrow[T c.Number] struct {
	Start Vector[T]
	End   Vector[T]

	HeadLength T

	HeadWidth T
}

func NewArrow[T c.Number](start, end Vector[T], headLength, headWidth T) Arrow[T] {
	return Arrow[T]{
		Start:      start,
		End:        end,
		HeadLength: headLength,
		HeadWidth:  headWidth,
	}
}
