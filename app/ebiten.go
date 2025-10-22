//go:build !headless

package app

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/core/ecs"

	"github.com/vistormu/xpeto/core/pkg/event"
	"github.com/vistormu/xpeto/core/pkg/window"
)

var toRunner = map[Runner]runner{
	Ebiten: &ebitenRunner{},
}

type ebitenRunner struct {
	app *App
}

func (r *ebitenRunner) Update() error {
	r.app.scheduler.RunUpdate(r.app.world)

	_, ok := event.GetEvents[EventExit](r.app.world)
	if ok {
		return ebiten.Termination
	}

	return nil
}

func (r *ebitenRunner) Draw(screen *ebiten.Image) {
	ecs.AddResource(r.app.world, screen)
	r.app.scheduler.RunDraw(r.app.world)
	ecs.RemoveResource[ebiten.Image](r.app.world)
}

func (r *ebitenRunner) Layout(w, h int) (int, int) {
	win, _ := ecs.GetResource[window.Window](r.app.world)
	return win.VWidth, win.VHeight
}

func (r *ebitenRunner) run(a *App) error {
	a.build()

	r.app = a

	return ebiten.RunGame(r)
}
