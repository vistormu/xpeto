package text

import (
	"image/color"

	"github.com/vistormu/xpeto/pkg/asset"
)

type Text struct {
	Font    asset.Handle
	Content string
	Color   color.Color
	Size    float64
}
