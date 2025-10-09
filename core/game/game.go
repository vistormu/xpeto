package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/event"
	"github.com/vistormu/xpeto/core/schedule"
)

// ====
// game
// ====
type Layout struct {
	Width  int
	Height int
}

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
	ecs.RemoveResource[ebiten.Image](g.world)
}

func (g *ebitenGame) Layout(w, h int) (int, int) {
	layout, _ := ecs.GetResource[Layout](g.world)
	return layout.Width, layout.Height
}

// =======
// builder
// =======
type Game struct {
	game *ebitenGame

	settings Settings
	plugins  []Plugin
}

func NewGame() *Game {
	return &Game{
		game:     nil,
		settings: Settings{},
		plugins:  make([]Plugin, 0),
	}
}

func (g *Game) WithPlugins(plugin ...Plugin) *Game {
	g.plugins = append(g.plugins, plugin...)
	return g
}

func (g *Game) WithSettings(settings Settings) *Game {
	g.settings = settings
	return g
}

func (g *Game) build() *Game {
	// game
	g.game = &ebitenGame{
		world:     ecs.NewWorld(),
		scheduler: schedule.NewScheduler(),
	}

	// core resources
	event.Startup(g.game.world)
	ecs.AddResource(g.game.world, &Layout{
		g.settings.VirtualWidth,
		g.settings.VirtualHeight,
	})

	// add event refreshing
	schedule.AddSystem(g.game.scheduler, schedule.Last, event.Update)

	// plugins
	for _, plugin := range g.plugins {
		plugin(g.game.world, g.game.scheduler)
	}

	g.game.scheduler.RunStartup(g.game.world)

	// settings
	ebiten.SetWindowSize(g.settings.WindowWidth, g.settings.WindowHeight)
	ebiten.SetWindowTitle(g.settings.WindowTitle)

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
