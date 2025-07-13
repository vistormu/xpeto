package collision

import (
	"github.com/vistormu/xpeto/internal/ecs"
	g "github.com/vistormu/xpeto/internal/geometry"
)

type CollisionEnter struct {
	Entity1     ecs.Entity
	Entity2     ecs.Entity
	Penetration g.Vector[float32]
	Normal      g.Vector[float32]
	Contact     g.Vector[float32]
	Static      bool
}

type CollisionStay struct {
	Entity1     ecs.Entity
	Entity2     ecs.Entity
	Penetration g.Vector[float32]
	Normal      g.Vector[float32]
	Contact     g.Vector[float32]
	Static      bool
}

type CollisionExit struct {
	Entity1 ecs.Entity
	Entity2 ecs.Entity
}
