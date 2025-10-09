package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"

	"github.com/vistormu/xpeto/core/pkg"
	"github.com/vistormu/xpeto/core/pkg/window"
)

// ====
// game
// ====
type ebitenGame struct {
	world     *ecs.World
	scheduler *schedule.Scheduler
}

func (g *ebitenGame) Update() error {
	g.scheduler.RunUpdate(g.world)

	return nil
}

func (g *ebitenGame) Draw(screen *ebiten.Image) {
	ecs.AddResource(g.world, screen)
	g.scheduler.RunDraw(g.world)
	ecs.RemoveResource[window.Screen](g.world)
}

func (g *ebitenGame) Layout(w, h int) (int, int) {
	layout, _ := ecs.GetResource[window.Layout](g.world)
	return layout.Width, layout.Height
}

// =======
// builder
// =======
type Game struct {
	game *ebitenGame
	pkgs []pkg.Pkg
}

func NewGame() *Game {
	return &Game{
		game: nil,
		pkgs: make([]pkg.Pkg, 0),
	}
}

func (g *Game) WithPkgs(pkgs ...pkg.Pkg) *Game {
	g.pkgs = append(g.pkgs, pkgs...)
	return g
}

func (g *Game) build() *Game {
	g.game = &ebitenGame{
		world:     ecs.NewWorld(),
		scheduler: schedule.NewScheduler(),
	}

	for _, pkg := range g.pkgs {
		pkg(g.game.world, g.game.scheduler)
	}

	g.game.scheduler.RunStartup(g.game.world)

	return g
}

func (g *Game) Run() error {
	g.build()

	err := ebiten.RunGame(g.game)
	if err != nil {
		return err
	}

	return nil
}
