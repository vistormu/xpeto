package physics

import (
	"image/color"

	// "github.com/hajimehoshi/ebiten/v2"
	"github.com/vistormu/xpeto/internal/core"
)

// Transient per-frame draw list (populated by the collector, consumed by the renderer)
type debugDrawList struct {
	Rects   []debugRect // AABB outlines
	Lines   []debugLine // contact normals, velocity vectors, etc.
	Strokes float32     // default stroke width
}

type debugRect struct {
	X, Y, W, H float32
	Col        color.Color
}

type debugLine struct {
	X1, Y1, X2, Y2 float32
	Col            color.Color
}

func debugRenderer(ctx *core.Context) {
	// screen := core.MustResource[*ebiten.Image](ctx)
}
