package game

import (
	"github.com/vistormu/xpeto"
)

type Fonts struct {
	Regular xp.Handle `path:"font.ttf"`
}

func loadASsets(w *xp.World) {
	xp.AddAssets[*xp.Font, Fonts](w)
}
