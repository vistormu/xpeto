// backend/window/apply.go
package window

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/window"
)

type lastRealWindow struct {
	window.RealWindow
}

func toEbitenResizingMode(m window.ResizingMode) ebiten.WindowResizingModeType {
	switch m {
	case window.ResizingModeDisabled:
		return ebiten.WindowResizingModeDisabled
	case window.ResizingModeOnlyFullscreenEnabled:
		return ebiten.WindowResizingModeOnlyFullscreenEnabled
	case window.ResizingModeEnabled:
		return ebiten.WindowResizingModeEnabled
	default:
		return ebiten.WindowResizingModeDisabled
	}
}

func applyInitial(w *ecs.World) {
	rw, _ := ecs.GetResource[window.RealWindow](w)

	ebiten.SetWindowSize(rw.Width, rw.Height)
	ebiten.SetWindowTitle(rw.Title)
	ebiten.SetFullscreen(rw.FullScreen)
	ebiten.SetVsyncEnabled(rw.VSync)
	ebiten.SetRunnableOnUnfocused(rw.RunnableOnUnfocused)

	ebiten.SetWindowResizingMode(toEbitenResizingMode(rw.ResizingMode))
	ebiten.SetWindowSizeLimits(rw.SizeLimits.MinW, rw.SizeLimits.MinH, rw.SizeLimits.MaxW, rw.SizeLimits.MaxH)

	// one-shot action
	applyAction(rw)

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

	if rw.ResizingMode != applied.ResizingMode {
		ebiten.SetWindowResizingMode(toEbitenResizingMode(rw.ResizingMode))
		applied.ResizingMode = rw.ResizingMode
	}

	if rw.SizeLimits != applied.SizeLimits {
		ebiten.SetWindowSizeLimits(rw.SizeLimits.MinW, rw.SizeLimits.MinH, rw.SizeLimits.MaxW, rw.SizeLimits.MaxH)
		applied.SizeLimits = rw.SizeLimits
	}

	// one-shot action (always check)
	if rw.Action != window.ActionNone {
		applyAction(rw)
	}
}

func applyAction(rw *window.RealWindow) {
	switch rw.Action {
	case window.ActionMaximize:
		ebiten.MaximizeWindow()
	case window.ActionMinimize:
		ebiten.MinimizeWindow()
	case window.ActionRestore:
		ebiten.RestoreWindow()
	}
	rw.Action = window.ActionNone
}
