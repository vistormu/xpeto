package ecs

import (
	"github.com/vistormu/xpeto/internal/core"
)

type Entity struct {
	Id      uint32
	Version uint32
}

type Component = any

type System = func(ctx *core.Context)
