package engine

import (
	"slices"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"
	"github.com/vistormu/xpeto/internal/event"
	"github.com/vistormu/xpeto/internal/scheduler"
)

// ====
// game
// ====
type ebitenGame struct {
	context       *core.Context
	scheduler     *scheduler.Scheduler
	drawScheduler *scheduler.Scheduler
}

func (g *ebitenGame) Update() error {
	g.scheduler.Run(g.context)

	return nil
}

func (g *ebitenGame) Draw(screen *ebiten.Image) {
	// g.drawScheduler.Run(g.context)
}

func (g *ebitenGame) Layout(w, h int) (int, int) {
	return 800, 600
}

// =======
// builder
// =======
type Game struct {
	game *ebitenGame

	settings Settings
	plugins  []core.Plugin
}

func NewGame() *Game {
	return &Game{
		game:     nil,
		settings: Settings{},
		plugins:  make([]core.Plugin, 0),
	}
}

func (g *Game) WithPlugins(plugin ...core.Plugin) *Game {
	g.plugins = append(g.plugins, plugin...)
	return g
}

func (g *Game) WithSettings(settings Settings) *Game {
	g.settings = settings
	return g
}

func (g *Game) build() *Game {
	g.game = new(ebitenGame)
	g.game.context = core.NewContext()
	g.game.scheduler = scheduler.NewScheduler(core.UpdateStages())
	startupScheduler := scheduler.NewScheduler(core.StartupStages())

	// core resources
	core.AddResource(g.game.context, ecs.NewWorld())
	core.AddResource(g.game.context, event.NewBus())

	// add packages
	for _, plugin := range g.plugins {
		sb := new(core.ScheduleBuilder)
		plugin.Build(g.game.context, sb)
		if slices.Contains(core.StartupStages(), sb.Stage) {
			startupScheduler.WithSchedule(&scheduler.Schedule{
				Name:      sb.Name,
				Stage:     sb.Stage,
				System:    sb.System,
				Before:    sb.Before,
				After:     sb.After,
				Condition: sb.Condition,
			})
		} else {
			g.game.scheduler.WithSchedule(&scheduler.Schedule{
				Name:      sb.Name,
				Stage:     sb.Stage,
				System:    sb.System,
				Before:    sb.Before,
				After:     sb.After,
				Condition: sb.Condition,
			})
		}
	}

	startupScheduler.Run(g.game.context)

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
