package game

import (
	"image/color"
	"math/rand"

	"github.com/vistormu/xpeto"
	"github.com/vistormu/xpeto/pkg/physics"
)

func createBackground(w *xp.World, r, g, b, a uint8) {
	win, _ := xp.GetResource[xp.Window](w)

	e := xp.AddEntity(w)
	xp.AddComponent(w, e, xp.Transform{
		X: float64(win.VWidth) / 2,
		Y: float64(win.VHeight) / 2,
	})
	xp.AddComponent(w, e, xp.Rect{
		Width:  float32(win.VWidth),
		Height: float32(win.VHeight),
		Fill:   color.RGBA{r, g, b, a},
		Layer:  0,
		Order:  0,
	})
}

func createWall(w *xp.World, x, y, wth, h float64) {
	win, _ := xp.GetResource[xp.Window](w)

	x *= float64(win.VWidth)
	y *= float64(win.VHeight)
	wth *= float64(win.VWidth)
	h *= float64(win.VHeight)

	e := xp.AddEntity(w)
	xp.AddComponent(w, e, xp.Transform{
		X: x,
		Y: y,
	})
	xp.AddComponent(w, e, physics.RigidBody{
		Type:        physics.Static,
		Restitution: 1,
	})
	xp.AddComponent(w, e, physics.Collider{
		Shape: physics.Rect{
			HalfW: wth / 2,
			HalfH: h / 2,
		},
	})
}

func createPlayer[T any](w *xp.World, x, y, wth, h float64) {
	win, _ := xp.GetResource[xp.Window](w)

	x *= float64(win.VWidth)
	y *= float64(win.VHeight)

	e := xp.AddEntity(w)
	xp.AddComponent(w, e, xp.Tag[T]{})
	xp.AddComponent(w, e, xp.Transform{
		X: x,
		Y: y,
	})
	xp.AddComponent(w, e, xp.Rect{
		Width:  float32(wth),
		Height: float32(h),
		Fill:   color.White,
		Layer:  1,
	})
	xp.AddComponent(w, e, physics.RigidBody{
		Type:        physics.Kinematic,
		Restitution: 1,
		Mass:        10,
	})
	xp.AddComponent(w, e, physics.Velocity{})
	xp.AddComponent(w, e, physics.Collider{
		Shape: physics.Rect{
			HalfW: wth / 2,
			HalfH: h / 2,
		},
	})
}

func createBall(w *xp.World, x, y float64) {
	win, _ := xp.GetResource[xp.Window](w)

	x *= float64(win.VWidth)
	y *= float64(win.VHeight)

	e := xp.AddEntity(w)
	xp.AddComponent(w, e, xp.Transform{
		X: x,
		Y: y,
	})
	xp.AddComponent(w, e, xp.Circle{
		Radius: 4,
		Fill:   color.White,
		Layer:  1,
	})
	xp.AddComponent(w, e, physics.RigidBody{
		Type:        physics.Dynamic,
		Mass:        1,
		Restitution: 1,
		Friction:    0,
	})
	xp.AddComponent(w, e, physics.Velocity{
		X: rand.NormFloat64()*3 + 100,
		Y: rand.Float64() * 20,
	})
	xp.AddComponent(w, e, physics.Collider{
		Shape: physics.Rect{
			HalfW: 2,
			HalfH: 2,
		},
	})
}

func createText(w *xp.World, text string, x, y float64) {
	fonts, ok := xp.GetResource[Fonts](w)
	if !ok {
		return
	}

	win, _ := xp.GetResource[xp.Window](w)

	x *= float64(win.VWidth)
	y *= float64(win.VHeight)

	e := xp.AddEntity(w)
	xp.AddComponent(w, e, xp.Transform{
		X: x,
		Y: y,
	})
	xp.AddComponent(w, e, xp.Text{
		Font:    fonts.Regular,
		Content: text,
		Align:   xp.AlignCenter,
		Color:   color.White,
		Size:    8,
		Layer:   1,
		Order:   0,
	})
}
