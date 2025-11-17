package geometry

import (
	c "github.com/vistormu/go-dsa/constraints"
)

type Path[T c.Number] struct {
	Points []Vector[T]
	Closed bool
}
