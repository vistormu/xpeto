package sprite

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"

	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/image"
	"github.com/vistormu/xpeto/pkg/render"
	"github.com/vistormu/xpeto/pkg/transform"
)

func packKey(layer, order uint16) uint64 {
	return (uint64(layer) << 16) | uint64(order)
}

func extractSprites(ctx *core.Context) []render.Renderable {
	w := core.MustResource[*ecs.World](ctx)
	entities := w.Query(ecs.And(
		ecs.Has[*Sprite](),
		ecs.Has[*transform.Transform](),
	))

	as, ok := core.GetResource[*asset.Server](ctx)
	if !ok {
		return nil
	}

	renderables := make([]render.Renderable, 0)
	for _, e := range entities {
		// components
		transform, _ := ecs.GetComponent[*transform.Transform](w, e)
		sprite, _ := ecs.GetComponent[*Sprite](w, e)

		// load asset
		img, ok := asset.GetAsset[*image.Image](as, sprite.Image)
		if !ok || img == nil {
			continue
		}

		item := &spriteItem{
			img:    img,
			x:      float64(transform.Position.X),
			y:      float64(transform.Position.Y),
			flipX:  sprite.FlipX,
			flipY:  sprite.FlipY,
			scaleX: float64(transform.Scale.X),
			scaleY: float64(transform.Scale.Y),
			layer:  sprite.Layer,
			order:  sprite.Order,
			key:    packKey(sprite.Layer, sprite.Order),
		}
		renderables = append(renderables, item)
	}

	return renderables
}
