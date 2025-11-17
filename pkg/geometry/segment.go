package geometry

import (
	c "github.com/vistormu/go-dsa/constraints"
)

type Segment[T c.Number] struct {
	Start Vector[T]
	End   Vector[T]
}
