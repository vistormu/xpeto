package game

import (
	"fmt"

	"github.com/vistormu/xpeto"
	"github.com/vistormu/xpeto/pkg/physics"
	"github.com/vistormu/xpeto/pkg/physics/debug"
)

// ======
// states
// ======
type gameState uint8

const (
	initial gameState = iota
	waiting
	playing
	scored
)

// ====
// tags
// ====
type player1 any
type player2 any

// =======
// startup
// =======
func initPhysiscs(w *xp.World) {
	sp, _ := xp.GetResource[physics.Space](w)
	// you would only need to change one of them
	sp.CellWidth = 16
	sp.CellHeight = 16
	sp.Height = 243
	sp.Width = 342

	// debug physics
	ds, _ := xp.GetResource[debug.Settings](w)
	ds.Enabled = true
}

func initWindow(w *xp.World) {
	win, _ := xp.GetResource[xp.Window](w)
	win.VWidth = 342
	win.VHeight = 243
}

func createInitialGameplayScene(w *xp.World) {
	// background
	createBackground(w, 40, 45, 52, 255)

	// collision borders
	createWall(w, 1, 0.5, 0.01, 1)
	createWall(w, 0, 0.5, 0.01, 1)

	// players
	createPlayer[player1](w, 0.1, 0.2, 5, 20)
	createPlayer[player2](w, 0.9, 0.8, 5, 20)

	// text
	createText(w,
		"press enter to start",
		0.5,
		0.5,
	)

	// TODO: create score text

	xp.SetNextState(w, waiting)
}

// ====
// loop
// ====
type EventScored[T any] struct {
	player T
}

func stateManager(w *xp.World) {
	current, _ := xp.GetState[gameState](w)

	if current == waiting {
		keyboard, _ := xp.GetResource[xp.Keyboard](w)
		if keyboard.IsJustPressed(xp.KeyEnter) {
			fmt.Println("enter pressed!")
			xp.SetNextState(w, playing)
		}
	}

	if current == playing {
		ev, ok := xp.GetEvents[EventScored[player1]](w)
		if ok {
			for _, e := range ev {
				// TODO
				fmt.Println(e)
				xp.SetNextState(w, scored)
			}
		}
	}
}

func onEnterPlaying(w *xp.World) {
	// TODO: remove text
	createBall(w, 0.5, 0.5)
}

func gameplay(w *xp.World) {
	keyboard, _ := xp.GetResource[xp.Keyboard](w)

	q := xp.Query1[physics.Velocity](w)
	for _, b := range q.Iter() {
		_, ok1 := xp.GetComponent[xp.Tag[player1]](w, b.Entity())
		_, ok2 := xp.GetComponent[xp.Tag[player2]](w, b.Entity())
		v := b.A()

		if keyboard.IsPressed(xp.KeyW) && ok1 {
			v.Y = -200
		} else if keyboard.IsPressed(xp.KeyS) && ok1 {
			v.Y = 200
		} else if ok1 {
			v.Y = 0
		}

		if keyboard.IsPressed(xp.KeyArrowUp) && ok2 {
			v.Y = -200
		} else if keyboard.IsPressed(xp.KeyArrowDown) && ok2 {
			v.Y = 200
		} else if ok2 {
			v.Y = 0
		}
	}
}
