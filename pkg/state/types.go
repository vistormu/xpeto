package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vistormu/xpeto/internal/core"
)

type FsmResource interface {
	Register(hook Hook, state any, fn any)
}

type FsmScheduler interface {
	OnEnter(ctx *core.Context)
	Update(ctx *core.Context)
	Draw(screen *ebiten.Image)
}
