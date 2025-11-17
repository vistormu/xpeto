package app

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"

	"github.com/vistormu/xpeto/core/pkg"
)

type App struct {
	world     *ecs.World
	scheduler *schedule.Scheduler
	runner    runner
	pkgs      []pkg.Pkg
}

func NewApp() *App {
	return &App{
		world:     ecs.NewWorld(),
		scheduler: schedule.NewScheduler(),
		pkgs:      make([]pkg.Pkg, 0),
	}
}

func (a *App) AddPkg(pkg pkg.Pkg) *App {
	a.pkgs = append(a.pkgs, pkg)
	return a
}

func (a *App) build() *App {
	// core packages
	pkg.CorePkgs(a.world, a.scheduler)

	// user packages
	for _, pkg := range a.pkgs {
		pkg(a.world, a.scheduler)
	}

	// startup schedules
	a.scheduler.RunStartup(a.world)

	return a
}

func (a *App) Run() error {
	defer a.scheduler.RunExit(a.world)

	err := a.runner.run(a)

	return err
}
