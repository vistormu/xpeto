//go:build !headless

package window

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func applyInitial(w *ecs.World) {
	rw, _ := ecs.GetResource[RealWindow](w)

	ebiten.SetWindowSize(rw.Width, rw.Height)
	ebiten.SetWindowTitle(rw.Title)
	ebiten.SetFullscreen(rw.FullScreen)

	ecs.AddResource(w, lastRealWindow{*rw})
}

func applyChanges(w *ecs.World) {
	rw, _ := ecs.GetResource[RealWindow](w)
	applied, _ := ecs.GetResource[lastRealWindow](w)

	if rw.Width != applied.Width || rw.Height != applied.Height {
		ebiten.SetWindowSize(rw.Width, rw.Height)
		applied.Width, applied.Height = rw.Width, rw.Height
	}

	if rw.Title != applied.Title {
		ebiten.SetWindowTitle(rw.Title)
		applied.Title = rw.Title
	}

	if rw.FullScreen != applied.FullScreen {
		ebiten.SetFullscreen(rw.FullScreen)
		applied.FullScreen = rw.FullScreen
	}
}

func setSystems(sch *schedule.Scheduler) {
	schedule.AddSystem(sch, schedule.PreStartup, applyInitial).Label("window.applyInitial")
	schedule.AddSystem(sch, schedule.First, applyChanges).Label("window.applyChanges")
}
