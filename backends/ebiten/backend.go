package ebiten

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/app"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func Backend() app.Backend {
	return &backend{}
}

type backend struct {
	w   *ecs.World
	sch *schedule.Scheduler
}

func (b *backend) Init(w *ecs.World, sch *schedule.Scheduler) {
	b.w = w
	b.sch = sch

	corePkgs(w, sch)
}

func (b *backend) Run() error {
	game := new(game)
	game.w = b.w
	game.sch = b.sch

	return ebiten.RunGame(game)
}
