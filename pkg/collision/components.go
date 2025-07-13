package collision

import (
	g "github.com/vistormu/xpeto/internal/geometry"
)

type Collidable struct {
	Size   g.Size[float32]
	Offset g.Vector[float32]
}
