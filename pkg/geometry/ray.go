package geometry

import (
	c "github.com/vistormu/go-dsa/constraints"
)

type Ray[T c.Number] struct {
	Origin    Vector[T]
	Direction Vector[T]
}
