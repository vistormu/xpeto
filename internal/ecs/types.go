package ecs

import (
	"github.com/vistormu/xpeto/internal/core"
)

type Entity = core.Handle[any]
type Component = any
type Archetype interface {
	Components() []Component
}
