//go:build !headless

package app

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/core/ecs"

	"github.com/vistormu/xpeto/core/pkg/event"
	"github.com/vistormu/xpeto/core/pkg/window"
)

type runner struct {
	app *App
}

func (r *runner) Update() error {
	r.app.scheduler.RunUpdate(r.app.world)

	_, ok := event.GetEvents[ExitApp](r.app.world)
	if ok {
		return ebiten.Termination
	}

	return nil
}

func (r *runner) Draw(screen *ebiten.Image) {
	ecs.AddResource(r.app.world, screen)
	r.app.scheduler.RunDraw(r.app.world)
	ecs.RemoveResource[ebiten.Image](r.app.world)
}

func (r *runner) Layout(w, h int) (int, int) {
	win, _ := ecs.GetResource[window.VirtualWindow](r.app.world)
	return win.Width, win.Height
}

func (r *runner) run(a *App) error {
	a.build()

	r.app = a

	return ebiten.RunGame(r)
}
