package font

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	ebitext "github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/transform"
)

type text struct {
	face    ebitext.Face
	content string
	align   Align
	color   color.Color
	x, y    float64
	layer   uint16
	order   uint16
}

func extractText(w *ecs.World) []text {
	q := ecs.NewQuery2[Text, transform.Transform](w)

	texts := make([]text, 0)
	for _, b := range q.Iter() {
		txt, tr := b.Components()

		// load asset
		fnt, ok := asset.GetAsset[Font](w, txt.Font)
		if !ok || fnt == nil {
			continue
		}

		texts = append(texts, text{
			face:    &ebitext.GoTextFace{Source: fnt.face, Size: txt.Size},
			content: txt.Content,
			align:   txt.Align,
			color:   txt.Color,
			x:       tr.X,
			y:       tr.Y,
			layer:   txt.Layer,
			order:   txt.Order,
		})
	}

	return texts
}

func sortText(t text) uint64 {
	return (uint64(t.layer) << 16) | uint64(t.order)
}

func drawText(screen *ebiten.Image, t text) {
	op := &ebitext.DrawOptions{}

	op.GeoM.Translate(t.x, t.y)
	op.ColorScale.ScaleWithColor(t.color)
	op.PrimaryAlign = t.align
	op.SecondaryAlign = t.align

	ebitext.Draw(screen, t.content, t.face, op)
}
