package geometry

import (
	c "github.com/vistormu/go-dsa/constraints"
)

type Capsule[T c.Number] struct {
	Segment Segment[T]
	Radius  T
}
