package ebiten

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/app"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"

	"github.com/vistormu/xpeto/core/event"
	"github.com/vistormu/xpeto/core/window"
)

// ===========
// ebiten game
// ===========
type game struct {
	w   *ecs.World
	sch *schedule.Scheduler
}

func (g *game) Update() error {
	g.sch.RunUpdate(g.w)

	_, ok := event.GetEvents[app.ExitApp](g.w)
	if ok {
		return ebiten.Termination
	}

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	ecs.AddResource(g.w, screen)
	g.sch.RunDraw(g.w)
	ecs.RemoveResource[ebiten.Image](g.w)
}

func (g *game) Layout(w, h int) (int, int) {
	win, _ := ecs.GetResource[window.VirtualWindow](g.w)
	return win.Width, win.Height
}

// =======
// backend
// =======
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
