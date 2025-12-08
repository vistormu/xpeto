package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/vistormu/go-dsa/set"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/pkg/input"
)

func updateKeyboardInput(w *ecs.World) {
	keyboard, _ := ecs.GetResource[input.Keyboard](w)

	keyboard.Clear()

	// rising edges
	var keys []ebiten.Key
	keys = inpututil.AppendJustPressedKeys(keys[:0])
	for _, k := range keys {
		keyboard.Press(input.Key(k))
	}

	// falling edges
	keys = inpututil.AppendJustReleasedKeys(keys[:0])
	for _, k := range keys {
		keyboard.Release(input.Key(k))
	}

	// keep pressed keys in sync if the window loses focus
	keys = inpututil.AppendPressedKeys(keys[:0])
	current := set.NewHashSet[input.Key]()
	for _, k := range keys {
		current.Add(input.Key(k))
		if !keyboard.IsPressed(input.Key(k)) {
			keyboard.Press(input.Key(k))
		}
	}

	// consider released any pressed key that is not in current
	for _, k := range keyboard.Pressed() {
		if !current.Contains(k) {
			keyboard.Release(k)
		}
	}

	// update durations
	for _, k := range keyboard.Pressed() {
		keyboard.SetDuration(k, inpututil.KeyPressDuration(ebiten.Key(k)))
	}
}
