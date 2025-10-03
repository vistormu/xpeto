package text

import (
	ebitext "github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"

	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/font"
	"github.com/vistormu/xpeto/pkg/render"
	"github.com/vistormu/xpeto/pkg/transform"
)

func packKey(layer, order uint16) uint64 {
	return (uint64(layer) << 16) | uint64(order)
}

func extractTexts(ctx *core.Context) []render.Renderable {
	w := core.MustResource[*ecs.World](ctx)
	entities := w.Query(ecs.And(
		ecs.Has[*Text](),
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
		text, _ := ecs.GetComponent[*Text](w, e)

		// load asset
		fnt, ok := asset.GetAsset[*font.Font](as, text.Font)
		if !ok || fnt == nil {
			continue
		}

		item := &textItem{
			face:    &ebitext.GoTextFace{Source: fnt.Face, Size: text.Size},
			content: text.Content,
			align:   text.Align,
			color:   text.Color,
			x:       float64(transform.Position.X),
			y:       float64(transform.Position.Y),
			layer:   text.Layer,
			order:   text.Order,
			key:     packKey(text.Layer, text.Order),
		}
		renderables = append(renderables, item)
	}

	return renderables
}
