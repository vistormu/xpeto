package physics

import (
	"github.com/vistormu/xpeto/core/pkg/transform"
)

// =====
// shape
// =====
type Shape interface {
	AABB(tr *transform.Transform) AABB
}

// ========
// collider
// ========
type Collider struct {
	Shape  Shape
	Layer  uint32
	Mask   uint32
	Sensor bool
}
