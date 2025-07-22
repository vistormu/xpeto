package physics

import (
	"github.com/vistormu/xpeto/internal/core"
)

type Dynamic struct {
	Velocity core.Vector[float32]
	Gravity  core.Vector[float32]
	Mass     float32
}
