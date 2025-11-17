package geometry

import (
	c "github.com/vistormu/go-dsa/constraints"
)

type Polygon[T c.Number] struct {
	Points []Vector[T]
}
