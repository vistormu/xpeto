package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"
	"github.com/vistormu/xpeto/internal/event"
	"github.com/vistormu/xpeto/internal/schedule"
)

// ====
// game
// ====
type Layout struct {
	Width  int
	Height int
}

type ebitenGame struct {
	context   *core.Context
	scheduler *schedule.Scheduler
}

func (g *ebitenGame) Update() error {
	g.scheduler.RunUpdate(g.context)

	return nil
}

func (g *ebitenGame) Draw(screen *ebiten.Image) {
	core.AddResource(g.context, screen)
	g.scheduler.RunDraw(g.context)
	core.RemoveResource[*ebiten.Image](g.context)
}

func (g *ebitenGame) Layout(w, h int) (int, int) {
	layout := core.MustResource[*Layout](g.context)
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
	g.game = new(ebitenGame)
	g.game.context = core.NewContext()
	g.game.scheduler = schedule.NewScheduler()

	// core resources
	core.AddResource(g.game.context, ecs.NewWorld())
	core.AddResource(g.game.context, event.NewBus())
	core.AddResource(g.game.context, &Layout{g.settings.VirtualWidth, g.settings.VirtualHeight})

	// plugins
	for _, plugin := range g.plugins {
		plugin(g.game.context, g.game.scheduler)
	}

	g.game.scheduler.RunStartup(g.game.context)

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
