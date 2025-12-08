package text

import (
	"image/color"

	"github.com/vistormu/xpeto/pkg/asset"
)

type Align uint8

const (
	AlignStart Align = iota
	AlignCenter
	AlignEnd
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
