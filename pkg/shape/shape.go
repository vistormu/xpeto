package shape

import (
	g "github.com/vistormu/xpeto/pkg/geometry"
)

type Shape struct {
	Shape  g.Geometry[float32]
	Fill   Fill
	Stroke Stroke
	Layer  uint16
	Order  uint16
}

// ===
// API
// ===
