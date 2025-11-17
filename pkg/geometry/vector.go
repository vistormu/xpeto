package geometry

import (
	c "github.com/vistormu/go-dsa/constraints"
)

type Vector[T c.Number] struct {
	X, Y T
}
