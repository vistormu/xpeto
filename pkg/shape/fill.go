package shape

import (
	"image/color"
)

type FillType uint8

const (
	FillSolid FillType = iota
	FillLinearGradient
	FillRadialGradient
	FillImage
	FillVideo
)

type Fill struct {
	Enabled bool
	Type    FillType
	Color   color.Color
	Opacity float32
}
