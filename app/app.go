package app

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"

	"github.com/vistormu/xpeto/core"
)

type App struct {
	world     *ecs.World
	scheduler *schedule.Scheduler
	backend   Backend
	pkgs      []core.Pkg
}

func NewApp(backend func() Backend) *App {
	return &App{
		world:     ecs.NewWorld(),
		scheduler: schedule.NewScheduler(),
		backend:   backend(),
		pkgs:      make([]core.Pkg, 0),
	}
}

func (a *App) AddPkg(pkg ...core.Pkg) *App {
	a.pkgs = append(a.pkgs, pkg...)
	return a
}

func (a *App) build() *App {
	// core packages
	core.CorePkgs(a.world, a.scheduler)

	// backend packages
	a.backend.Init(a.world, a.scheduler)

	// user packages
	for _, pkg := range a.pkgs {
		pkg(a.world, a.scheduler)
	}

	// startup schedules
	a.scheduler.RunStartup(a.world)

	return a
}

func (a *App) Run() error {
	a.build()
	defer a.scheduler.RunExit(a.world)

	return a.backend.Run()
}
