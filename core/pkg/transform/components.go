package transform

import (
	"github.com/vistormu/xpeto/internal/core"
)

type Transform struct {
	Position core.Vector[float32]
	Scale    core.Vector[float32]
	Rotation float32
}

// TODO: implement this correctly
type GlobalTransform struct {
	Position core.Vector[float32]
	Scale    core.Vector[float32]
	Rotation float32
}
