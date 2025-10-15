package app

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/core/ecs"

	"github.com/vistormu/xpeto/core/pkg/window"
)

type ebitenRunner struct {
	app *App
}

func (r *ebitenRunner) Update() error {
	r.app.scheduler.RunUpdate(r.app.world)
	return nil
}

func (r *ebitenRunner) Draw(screen *ebiten.Image) {
	ecs.AddResource(r.app.world, screen)
	r.app.scheduler.RunDraw(r.app.world)
	ecs.RemoveResource[window.Screen](r.app.world)
}

func (r *ebitenRunner) Layout(w, h int) (int, int) {
	layout, _ := ecs.GetResource[window.Layout](r.app.world)
	return layout.Width, layout.Height
}

func (r *ebitenRunner) run(a *App) error {
	a.build()

	r.app = a

	return ebiten.RunGame(r)
}
