package sprite

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/pkg/transform"

	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/image"
	"github.com/vistormu/xpeto/pkg/render"
)

func packKey(layer, order uint16) uint64 {
	return (uint64(layer) << 16) | uint64(order)
}

func extractSprites(w *ecs.World) []render.Renderable {
	q := ecs.NewQuery2[Sprite, transform.Transform](w)

	renderables := make([]render.Renderable, 0)
	for _, b := range q.Iter() {
		sprite := b.A()
		tr := b.B()

		img, ok := asset.GetAsset[*image.Image](w, sprite.Image)
		if !ok || img == nil {
			continue
		}

		item := &spriteItem{
			img:   img,
			x:     tr.X,
			y:     tr.Y,
			flipX: sprite.FlipX,
			flipY: sprite.FlipY,
			layer: sprite.Layer,
			order: sprite.Order,
			key:   packKey(sprite.Layer, sprite.Order),
		}
		renderables = append(renderables, item)
	}

	return renderables
}
