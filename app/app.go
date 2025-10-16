package app

import (
	"errors"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"

	"github.com/vistormu/xpeto/core/pkg"
)

type App struct {
	world     *ecs.World
	scheduler *schedule.Scheduler
	runner    Runner
	pkgs      []pkg.Pkg
}

func NewApp() *App {
	return &App{
		world:     ecs.NewWorld(),
		scheduler: schedule.NewScheduler(),
		pkgs:      make([]pkg.Pkg, 0),
	}
}

func (a *App) WithPkgs(pkgs ...pkg.Pkg) *App {
	a.pkgs = append(a.pkgs, pkgs...)
	return a
}

func (a *App) WithRunner(runner Runner) *App {
	a.runner = runner
	return a
}

func (a *App) build() *App {
	// core packages
	for _, pkg := range pkg.CorePkgs() {
		pkg(a.world, a.scheduler)
	}

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
	runner, ok := toRunner[a.runner]
	if !ok {
		return errors.New("runner not found")
	}

	err := runner.run(a)

	return err
}
