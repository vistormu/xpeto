package font

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/pkg/asset"
	xptext "github.com/vistormu/xpeto/pkg/text"
	"github.com/vistormu/xpeto/pkg/transform"
)

type renderable struct {
	face    text.Face
	content string
	align   xptext.Align
	color   color.Color
	x, y    float64
	layer   uint16
	order   uint16
}

var xpToEbiAlign = map[xptext.Align]text.Align{
	xptext.AlignStart:  text.AlignStart,
	xptext.AlignCenter: text.AlignStart,
	xptext.AlignEnd:    text.AlignEnd,
}

func extractText(w *ecs.World) []renderable {
	q := ecs.NewQuery2[xptext.Text, transform.Transform](w)

	texts := make([]renderable, 0)
	for _, b := range q.Iter() {
		txt, tr := b.Components()

		// load asset
		fnt, ok := asset.GetAsset[font](w, txt.Font)
		if !ok || fnt == nil {
			continue
		}

		texts = append(texts, renderable{
			face:    &text.GoTextFace{Source: fnt.GoTextFaceSource, Size: txt.Size},
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

func sortText(t renderable) uint64 {
	return (uint64(t.layer) << 16) | uint64(t.order)
}

func drawText(screen *ebiten.Image, t renderable) {
	op := &text.DrawOptions{}

	op.GeoM.Translate(t.x, t.y)
	op.ColorScale.ScaleWithColor(t.color)
	op.PrimaryAlign = xpToEbiAlign[t.align]
	op.SecondaryAlign = xpToEbiAlign[t.align]

	text.Draw(screen, t.content, t.face, op)
}
