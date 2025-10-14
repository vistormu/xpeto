package window

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vistormu/xpeto/core/ecs"
)

type WindowSettings struct {
	Title        string
	Width        int
	Height       int
	VWidth       int
	VHeight      int
	FullScreen   bool
	AntiAliasing bool
}

type lastSettings struct {
	WindowSettings
}

func applyInitialSettings(w *ecs.World) {
	ws, _ := ecs.GetResource[WindowSettings](w)

	ebiten.SetWindowSize(ws.Width, ws.Height)
	ebiten.SetWindowTitle(ws.Title)
	ebiten.SetFullscreen(ws.FullScreen)

	ecs.AddResource(w, lastSettings{*ws})
}

func applyChanges(w *ecs.World) {
	ws, _ := ecs.GetResource[WindowSettings](w)
	applied, _ := ecs.GetResource[lastSettings](w)

	if ws.Width != applied.Width || ws.Height != applied.Height {
		ebiten.SetWindowSize(ws.Width, ws.Height)
		applied.Width, applied.Height = ws.Width, ws.Height
	}

	if ws.Title != applied.Title {
		ebiten.SetWindowTitle(ws.Title)
		applied.Title = ws.Title
	}

	if ws.FullScreen != applied.FullScreen {
		ebiten.SetFullscreen(ws.FullScreen)
		applied.FullScreen = ws.FullScreen
	}
}
