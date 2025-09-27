package ecs

import (
	"github.com/vistormu/xpeto/internal/core"
)

type Entity = core.Id
type Component = any
type System = func(ctx *core.Context)
