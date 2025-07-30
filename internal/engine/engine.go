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
	g.drawScheduler.Run(g.context)
}

func (g *ebitenGame) Layout(w, h int) (int, int) {
	return 800, 600
}

// =======
// builder
// =======
type Game struct {
	resources []any
	settings  Settings
	schedules []*scheduler.Schedule
	pkgs      []Pkg
}

func NewGame() *Game {
	return &Game{
		resources: make([]any, 0),
		settings:  Settings{},
		schedules: make([]*scheduler.Schedule, 0),
		pkgs:      make([]Pkg, 0),
	}
}

func (g *Game) WithResources(resources ...any) *Game {
	g.resources = append(g.resources, resources...)
	return g
}

func (g *Game) WithSystem(name string, stage scheduler.Stage, run ecs.System, condition func(*core.Context) bool) *Game {
	g.schedules = append(g.schedules, &scheduler.Schedule{
		Name:      name,
		Stage:     stage,
		System:    run,
		Before:    nil, // no before systems
		After:     nil, // no after systems
		Condition: condition,
	})

	return g
}

func (g *Game) WithPkgs(pkg ...Pkg) *Game {
	g.pkgs = append(g.pkgs, pkg...)
	return g
}

func (g *Game) WithSettings(settings Settings) *Game {
	g.settings = settings
	return g
}

func (g *Game) Run() error {
	game := new(ebitenGame)
	game.context = core.NewContext()
	game.scheduler = scheduler.NewScheduler(scheduler.UpdateStages())
	startupScheduler := scheduler.NewScheduler(scheduler.StartupStages())

	// core resources
	core.AddResource(game.context, ecs.NewWorld())
	core.AddResource(game.context, event.NewBus())

	// user resources
	for _, res := range g.resources {
		core.AddResource(game.context, res)
	}

	// add schedules
	for _, sch := range g.schedules {
		if slices.Contains(scheduler.StartupStages(), sch.Stage) {
			startupScheduler.WithSchedule(sch)
		} else {
			game.scheduler.WithSchedule(sch)
		}
	}

	// add packages
	for _, pkg := range g.pkgs {
		g.resources = append(g.resources, pkg.Resources()...)

		for _, sch := range pkg.Schedules() {
			if slices.Contains(scheduler.StartupStages(), sch.Stage) {
				startupScheduler.WithSchedule(sch)
			} else {
				game.scheduler.WithSchedule(sch)
			}
		}

		pkg.Build(game.context)
	}

	// settings
	ebiten.SetWindowSize(g.settings.WindowWidth, g.settings.WindowHeight)
	ebiten.SetWindowTitle(g.settings.WindowTitle)

	err := ebiten.RunGame(game)
	if err != nil {
		return err
	}

	return nil
}
