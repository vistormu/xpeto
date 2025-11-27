package shape

import (
	"image/color"
)

// =====
// shape
// =====
type Shape struct {
	Visible bool
	Fill    []Fill
	Stroke  []Stroke
	Layer   uint16
	Order   uint16
}

func newShape() Shape {
	return Shape{
		Visible: true,
		Fill:    make([]Fill, 0),
		Stroke:  make([]Stroke, 0),
	}
}

func (s Shape) AddFillSolid(c color.Color) Shape {
	s.Fill = append(s.Fill, Fill{
		Visible: true,
		Type:    FillSolid,
		Color:   c,
	})

	return s
}

func (s Shape) AddStroke(c color.Color, w float32) Shape {
	s.Stroke = append(s.Stroke, Stroke{
		Visible: true,
		Color:   c,
		Width:   w,
	})

	return s
}

func (s Shape) AddOrder(l uint16, o uint16) Shape {
	s.Layer = l
	s.Order = o

	return s
}

// ====
// fill
// ====
type FillType uint8

const (
	FillSolid FillType = iota
	FillLinearGradient
	FillRadialGradient
	FillImage
	FillVideo
)

type Fill struct {
	Visible bool
	Type    FillType
	Color   color.Color
}

// ======
// stroke
// ======
type Stroke struct {
	Visible bool
	Color   color.Color
	Width   float32
}
