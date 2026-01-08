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
	Type  FillType
	Color color.Color
}
