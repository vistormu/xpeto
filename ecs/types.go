package ecs

import (
	"github.com/vistormu/xpeto/image"
	"github.com/vistormu/xpeto/internal/core"
)

type Entity = core.Handle[any]

type System interface {
	OnLoad(context *Context)
	OnUnload(context *Context)

	Draw(screen *image.Image)

	Update(context *Context, dt float32)
	FixedUpdate(context *Context, dt float32)
}
