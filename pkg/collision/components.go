package collision

import (
	"github.com/vistormu/xpeto/internal/core"
)

type Collidable struct {
	Size   core.Size[float32]
	Offset core.Vector[float32]
}
