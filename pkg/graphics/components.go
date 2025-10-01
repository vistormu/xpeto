package graphics

import "image/color"

type Circle struct {
	Radius    int
	Fill      color.Color
	Stroke    color.Color
	Linewidth int
}

type Rect struct {
	Width  float32
	Height float32
}
