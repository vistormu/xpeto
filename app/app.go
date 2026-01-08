package app

import (
	"fmt"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/log"
	"github.com/vistormu/xpeto/core/schedule"

	"github.com/vistormu/xpeto/core"
)

// =======
// options
// =======
type option = func(*App)

type appOpt struct{}

var AppOpt appOpt

func (appOpt) Pkgs(pkgs ...core.Pkg) option {
	return func(a *App) {
		a.pkgs = append(a.pkgs, pkgs...)
	}
}

// ===
// app
// ===
type App struct {
	backend BackendFactory
	pkgs    []core.Pkg
}

func NewApp(backend BackendFactory, opts ...option) *App {
	app := &App{
		backend: backend,
		pkgs:    make([]core.Pkg, 0),
	}

	for _, opt := range opts {
		if opt != nil {
			opt(app)
		}
	}

	return app
}

func (a *App) Run() error {
	// core dependencies
	w := ecs.NewWorld()
	sch := schedule.NewScheduler()

	core.CorePkgs(w, sch)

	// backend
	if a.backend == nil {
		return fmt.Errorf("nil backend\n")
	}
	backend, err := a.backend(w, sch)
	if err != nil {
		return err
	}

	// optional dependencies
	for _, pkg := range a.pkgs {
		if pkg != nil {
			pkg(w, sch)
		}
	}

	// schedules
	if ds := schedule.Diagnostics(sch); len(ds) != 0 {
		msg := "errors detected during schedule compilation\n"
		for _, d := range ds {
			msg += fmt.Sprintf("-> %s\n", d.Message)
			msg += fmt.Sprintf("   |> system: %s (id: %d)\n", d.Label, d.Id)
			msg += fmt.Sprintf("   |> stage: %s\n\n", d.Stage.String())
		}
		return fmt.Errorf(msg)
	}

	defer schedule.RunExit(w, sch)
	defer log.RecoverLog(w)()

	schedule.RunStartup(w, sch)

	return backend.Run()
}
