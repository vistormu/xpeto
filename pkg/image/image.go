package image

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/pkg/asset"
)

// =====
// image
// =====
type Image struct {
	Img *ebiten.Image
}

// ======
// sprite
// ======
type Sprite struct {
	Image asset.Asset
	FlipX bool
	FlipY bool
	Layer uint16
	Order uint16
}
