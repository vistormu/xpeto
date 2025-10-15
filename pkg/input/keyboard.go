package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/vistormu/go-dsa/set"

	"github.com/vistormu/xpeto/core/ecs"
	// "github.com/vistormu/xpeto/internal/event"
)

type Key = ebiten.Key

const (
	KeyA = ebiten.KeyA
	KeyD = ebiten.KeyD
	KeyS = ebiten.KeyS
	KeyW = ebiten.KeyW

	KeyArrowDown = ebiten.KeyArrowDown
	KeyArrowUp   = ebiten.KeyArrowUp
)

type Keyboard = ButtonInput[Key]

func newKeyboard() *Keyboard {
	return newButtonInput[Key]()
}

func updateKeyboardInput(w *ecs.World) {
	keyboard, _ := ecs.GetResource[Keyboard](w)

	keyboard.clear()

	// rising edges
	var keys []Key
	keys = inpututil.AppendJustPressedKeys(keys[:0])
	for _, k := range keys {
		keyboard.press(k)
		// event.Publish(eb, KeyJustPressed{Key: k})
	}

	// falling edges
	keys = inpututil.AppendJustReleasedKeys(keys[:0])
	for _, k := range keys {
		keyboard.release(k)
		// event.Publish(eb, KeyJustReleased{Key: k})
	}

	// keep pressed keys in sync if the window loses focus
	keys = inpututil.AppendPressedKeys(keys[:0])
	current := set.NewHashSet[Key]()
	for _, k := range keys {
		current.Add(k)
		if !keyboard.IsPressed(k) {
			keyboard.press(k)
		}
	}

	// consider released any pressed key that is not in current
	for _, k := range keyboard.pressed.Values() {
		if !current.Contains(k) {
			keyboard.release(k)
		}
	}

	// update durations
	for _, k := range keyboard.pressed.Values() {
		keyboard.setDuration(k, inpututil.KeyPressDuration(k))
	}
}
