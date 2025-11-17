package geometry

import (
	c "github.com/vistormu/go-dsa/constraints"
)

type Line[T c.Number] struct {
	Point     Vector[T]
	Direction Vector[T]
}
