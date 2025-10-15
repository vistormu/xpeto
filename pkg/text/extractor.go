package text

import (
	ebitext "github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/pkg/transform"

	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/font"
	"github.com/vistormu/xpeto/pkg/render"
)

func packKey(layer, order uint16) uint64 {
	return (uint64(layer) << 16) | uint64(order)
}

func extractTexts(w *ecs.World) []render.Renderable {
	q := ecs.NewQuery2[Text, transform.Transform](w)

	renderables := make([]render.Renderable, 0)
	for _, b := range q.Iter() {
		// components
		txt := b.A()
		tr := b.B()

		// load asset
		fnt, ok := asset.GetAsset[*font.Font](w, txt.Font)
		if !ok || fnt == nil {
			continue
		}

		item := &textItem{
			face:    &ebitext.GoTextFace{Source: fnt.Face, Size: txt.Size},
			content: txt.Content,
			align:   txt.Align,
			color:   txt.Color,
			x:       tr.X,
			y:       tr.Y,
			layer:   txt.Layer,
			order:   txt.Order,
			key:     packKey(txt.Layer, txt.Order),
		}
		renderables = append(renderables, item)
	}

	return renderables
}
