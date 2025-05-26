package image

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vistormu/xpeto/internal/core"
)

type Image = ebiten.Image
type Handle = core.Handle[*Image]
