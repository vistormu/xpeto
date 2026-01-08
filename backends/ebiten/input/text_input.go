package input

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/event"
	"github.com/vistormu/xpeto/pkg/input"
)

func emitText(w *ecs.World) {
	rs := ebiten.AppendInputChars(nil)
	for _, r := range rs {
		event.AddEvent(w, input.TextInputEvent{Rune: r})
	}
}
