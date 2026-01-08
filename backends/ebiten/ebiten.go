package ebiten

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/app"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"

	"github.com/vistormu/xpeto/core/event"
	"github.com/vistormu/xpeto/core/window"
)

// =======
// backend
// =======
func Backend(w *ecs.World, sch *schedule.Scheduler) (app.Backend, error) {
	game := &game{
		w:   w,
		sch: sch,
	}

	corePkgs(game.w, game.sch)

	return game, nil
}

// ======
// screen
// ======
type screenBuffer struct {
	screen *ebiten.Image
	w, h   int
}

// ====
// game
// ====
type game struct {
	w   *ecs.World
	sch *schedule.Scheduler
}

func (g *game) Update() error {
	schedule.RunUpdate(g.w, g.sch)

	_, ok := event.GetEvents[app.ExitAppEvent](g.w)
	if ok {
		return ebiten.Termination
	}

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	// update buffer
	vw, vh := window.GetVirtualWindowSize[int](g.w)
	sb, _ := ecs.GetResource[screenBuffer](g.w)
	if sb.screen == nil || sb.w != vw || sb.h != vh {
		sb.screen = ebiten.NewImage(vw, vh)
		sb.w = vw
		sb.h = vh
	}

	// draw
	sb.screen.Clear()
	ecs.AddResource(g.w, sb.screen)
	schedule.RunDraw(g.w, g.sch)
	ecs.RemoveResource[ebiten.Image](g.w)

	// update viewport
	vp, _ := ecs.GetResource[window.Viewport](g.w)

	op := &ebiten.DrawImageOptions{}
	if vp.Scale > 0 {
		op.GeoM.Scale(vp.ScaleF, vp.ScaleF)
	}
	op.GeoM.Translate(vp.OffsetX, vp.OffsetY)

	screen.Clear()
	screen.DrawImage(sb.screen, op)
}

func (g *game) Layout(outsideW, outsideH int) (int, int) {
	obs, _ := ecs.GetResource[window.RealWindowObserved](g.w)
	obs.Width = outsideW
	obs.Height = outsideH
	obs.DeviceScale = ebiten.Monitor().DeviceScaleFactor()

	vw, vh, ok := window.GetDesiredVirtualSize(g.w)
	if ok {
		window.SetVirtualWindowSize(g.w, vw, vh)
	}

	vp := window.ComputeViewport(g.w)
	ecs.AddResource(g.w, vp)

	// w, h := window.GetVirtualWindowSize[int](g.w)

	return outsideW, outsideH
}

func (g *game) Run() error {
	return ebiten.RunGame(g)
}
