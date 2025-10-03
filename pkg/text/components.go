package text

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	ebitext "github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/vistormu/xpeto/pkg/asset"
)

type Align = ebitext.Align

var AlignStart Align = ebitext.AlignStart
var AlignCenter Align = ebitext.AlignCenter
var AlignEnd Align = ebitext.AlignEnd

type Text struct {
	Font    asset.Handle
	Content string
	Align   Align
	Color   color.Color
	Size    float64
	Layer   uint16
	Order   uint16
}

type textItem struct {
	face    ebitext.Face
	content string
	align   Align
	color   color.Color
	x, y    float64
	layer   uint16
	order   uint16
	key     uint64
}

func (ti *textItem) Draw(screen *ebiten.Image) {
	op := &ebitext.DrawOptions{}

	op.GeoM.Translate(ti.x, ti.y)
	op.ColorScale.ScaleWithColor(ti.color)
	op.PrimaryAlign = ti.align
	op.SecondaryAlign = ti.align

	ebitext.Draw(screen, ti.content, ti.face, op)
}

func (ti *textItem) SortKey() uint64 {
	return ti.key
}
