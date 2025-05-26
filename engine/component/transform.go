package component

import (
	g "github.com/vistormu/xpeto/internal/geometry"
)

type Transform struct {
	Position g.Vector[float32]
	Scale    g.Vector[float32]
	Rotation float32
}
