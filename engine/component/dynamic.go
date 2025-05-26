package component

import (
	g "github.com/vistormu/xpeto/internal/geometry"
)

type Dynamic struct {
	Velocity g.Vector[float32]
	Gravity  g.Vector[float32]
	Mass     float32
}
