package sprite

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/image"
)

type Sprite struct {
	Image asset.Handle
	FlipX bool
	FlipY bool
	Layer uint16
	Order uint16
}

type spriteItem struct {
	img    *image.Image
	x, y   float64
	flipX  bool
	flipY  bool
	scaleX float64
	scaleY float64
	layer  uint16
	order  uint16
	key    uint64
}

func (si *spriteItem) Draw(screen *ebiten.Image) {
	w := si.img.Img.Bounds().Dx()
	h := si.img.Img.Bounds().Dy()

	var op ebiten.DrawImageOptions

	sx := si.scaleX
	sy := si.scaleY
	if si.flipX {
		sx = -sx
	}
	if si.flipY {
		sy = -sy
	}

	op.GeoM.Scale(sx, sy)

	if si.flipX {
		op.GeoM.Translate(float64(w)*si.scaleX, 0)
	}
	if si.flipY {
		op.GeoM.Translate(0, float64(h)*si.scaleY)
	}

	op.GeoM.Translate(si.x, si.y)

	screen.DrawImage(si.img.Img, &op)
}

func (si *spriteItem) SortKey() uint64 {
	return si.key
}
