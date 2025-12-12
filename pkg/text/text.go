package text

import (
	"image/color"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/render"
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
	Anchor  render.Anchor
}

type DefaultFonts struct {
	Default asset.Asset `path:"default/font.ttf"`
	Pico8   asset.Asset `path:"default/pico8.ttf"`
}

func loadFonts(w *ecs.World) {
	asset.AddAsset[DefaultFonts](w)
}
