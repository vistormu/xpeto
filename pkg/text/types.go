package text

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Font struct {
	Font    *opentype.Font
	Face    font.Face
	Size    float64
	DPI     float64
	Hinting font.Hinting
}
