package collision

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"
)

type CollisionEnter struct {
	Entity1     ecs.Entity
	Entity2     ecs.Entity
	Penetration core.Vector[float32]
	Normal      core.Vector[float32]
	Contact     core.Vector[float32]
	Static      bool
}

type CollisionStay struct {
	Entity1     ecs.Entity
	Entity2     ecs.Entity
	Penetration core.Vector[float32]
	Normal      core.Vector[float32]
	Contact     core.Vector[float32]
	Static      bool
}

type CollisionExit struct {
	Entity1 ecs.Entity
	Entity2 ecs.Entity
}
