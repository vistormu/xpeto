package window

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/window"
)

type lastRealWindow struct {
	window.RealWindow
}

func applyInitial(w *ecs.World) {
	rw, _ := ecs.GetResource[window.RealWindow](w)

	ebiten.SetWindowSize(rw.Width, rw.Height)
	ebiten.SetWindowTitle(rw.Title)
	ebiten.SetFullscreen(rw.FullScreen)
	ebiten.SetVsyncEnabled(rw.VSync)
	ebiten.SetRunnableOnUnfocused(rw.RunnableOnUnfocused)

	ecs.AddResource(w, lastRealWindow{*rw})
}

func applyChanges(w *ecs.World) {
	rw, _ := ecs.GetResource[window.RealWindow](w)
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

	if rw.VSync != applied.VSync {
		ebiten.SetVsyncEnabled(rw.VSync)
		applied.VSync = rw.VSync
	}

	if rw.RunnableOnUnfocused != applied.RunnableOnUnfocused {
		ebiten.SetRunnableOnUnfocused(rw.RunnableOnUnfocused)
		applied.RunnableOnUnfocused = rw.RunnableOnUnfocused
	}
}
