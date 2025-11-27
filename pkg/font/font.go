package font

import (
	"image/color"

	ebitext "github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/vistormu/xpeto/pkg/asset"
)

// ====
// font
// ====
type Font struct {
	face *ebitext.GoTextFaceSource
}

// ====
// text
// ====
type Align = ebitext.Align

const (
	AlignStart  Align = ebitext.AlignStart
	AlignCenter Align = ebitext.AlignCenter
	AlignEnd    Align = ebitext.AlignEnd
)

type Text struct {
	Font    asset.Asset
	Content string
	Align   Align
	Color   color.Color
	Size    float64
	Layer   uint16
	Order   uint16
}
