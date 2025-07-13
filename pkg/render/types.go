package render

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vistormu/xpeto/internal/core"
)

type Renderer = ebiten.Image
type Image = core.Handle[*Renderer]
