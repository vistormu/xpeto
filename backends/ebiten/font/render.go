package font

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/vistormu/xpeto/backends/ebiten/shared"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/log"
	"github.com/vistormu/xpeto/core/window"
	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/render"
	xptext "github.com/vistormu/xpeto/pkg/text"
	"github.com/vistormu/xpeto/pkg/transform"
)

type renderable struct {
	face    text.Face
	content string
	align   xptext.Align
	anchor  render.Anchor
	color   color.Color
	x, y    float64
	layer   uint16
	order   uint16
	snap    bool
}

var xpToEbiAlign = map[xptext.Align]text.Align{
	xptext.AlignStart:  text.AlignStart,
	xptext.AlignCenter: text.AlignStart,
	xptext.AlignEnd:    text.AlignEnd,
}

func extractText(w *ecs.World) []renderable {
	sc, _ := ecs.GetResource[window.Scaling](w)
	q := ecs.NewQuery2[xptext.Text, transform.Transform](w)

	texts := make([]renderable, 0)
	for _, b := range q.Iter() {
		txt, tr := b.Components()

		fnt, ok := asset.GetAsset[font](w, txt.Font)
		if !ok || fnt == nil {
			log.LogError(w, "tried to load a missing font", log.F("function", "backends/ebiten/font/render.go:extractText"))
			continue
		}

		face := fnt.Face(txt.Size)
		texts = append(texts, renderable{
			face:    face,
			content: txt.Content,
			align:   txt.Align,
			anchor:  txt.Anchor,
			color:   txt.Color,
			x:       tr.X,
			y:       tr.Y,
			layer:   txt.Layer,
			order:   txt.Order,
			snap:    sc.SnapPixels,
		})
	}

	return texts
}

func sortText(t renderable) uint64 {
	return (uint64(t.layer) << 16) | uint64(t.order)
}

func drawText(screen *ebiten.Image, t renderable) {
	op := &text.DrawOptions{}

	op.ColorScale.ScaleWithColor(t.color)
	op.PrimaryAlign = xpToEbiAlign[t.align]
	op.SecondaryAlign = xpToEbiAlign[t.align]

	w, h := text.Measure(t.content, t.face, 0)
	dx, dy := shared.Offset(w, h, t.anchor)

	x := t.x + dx
	y := t.y + dy

	if t.snap {
		x = math.Round(x)
		y = math.Round(y)
	}

	op.GeoM.Translate(x, y)
	text.Draw(screen, t.content, t.face, op)
}
