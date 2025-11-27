package image

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/transform"
)

type sprite struct {
	img   *Image
	x, y  float64
	flipX bool
	flipY bool
	// scaleX float64
	// scaleY float64
	layer uint16
	order uint16
}

func extractSprite(w *ecs.World) []sprite {
	q := ecs.NewQuery2[Sprite, transform.Transform](w)

	sprites := make([]sprite, 0)
	for _, b := range q.Iter() {
		s, tr := b.Components()

		img, ok := asset.GetAsset[Image](w, s.Image)
		if !ok || img == nil {
			continue
		}

		sprites = append(sprites, sprite{
			img:   img,
			x:     tr.X,
			y:     tr.Y,
			flipX: s.FlipX,
			flipY: s.FlipY,
			layer: s.Layer,
			order: s.Order,
		})
	}

	return sprites
}

func sortSprite(s sprite) uint64 {
	return (uint64(s.layer) << 16) | uint64(s.order)
}

func drawSprite(screen *ebiten.Image, s sprite) {
	// w := si.img.Img.Bounds().Dx()
	// h := si.img.Img.Bounds().Dy()

	var op ebiten.DrawImageOptions

	// sx := si.scaleX
	// sy := si.scaleY
	// if si.flipX {
	// 	sx = -sx
	// }
	// if si.flipY {
	// 	sy = -sy
	// }

	// op.GeoM.Scale(sx, sy)

	// if si.flipX {
	// 	op.GeoM.Translate(float64(w)*si.scaleX, 0)
	// }
	// if si.flipY {
	// 	op.GeoM.Translate(0, float64(h)*si.scaleY)
	// }

	op.GeoM.Translate(s.x, s.y)

	screen.DrawImage(s.img.Img, &op)
}
