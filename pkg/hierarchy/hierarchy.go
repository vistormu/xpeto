package hierarchy

import (
	"github.com/vistormu/xpeto/core/ecs"
)

type ChildOf struct {
	Entity ecs.Entity
}

type Children struct {
	Entities []ecs.Entity
}
