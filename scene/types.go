package scene

import (
	"github.com/vistormu/xpeto/ecs"
	"github.com/vistormu/xpeto/image"
)

type Scene interface {
	OnLoad(ctx *ecs.Context)
	OnUnload(ctx *ecs.Context)
	OnEnter(ctx *ecs.Context)
	OnExit(ctx *ecs.Context)
	Update(ctx *ecs.Context, dt float32)
	Draw(screen *image.Image)
	Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int)
}
