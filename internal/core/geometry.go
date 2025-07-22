package core

import (
	"github.com/vistormu/go-dsa/constraints"
	g "github.com/vistormu/go-dsa/geometry"
)

// vector
type Vector[T constraints.Number] = g.Vector2D[T]

func NewVector[T constraints.Number](x, y T) Vector[T] {
	return g.NewVector2D(x, y)
}

// size
type Size[T constraints.Number] = g.Size2D[T]

func NewSize[T constraints.Number](width, height T) Size[T] {
	return g.NewSize2D(width, height)
}

// rect
type Rect[T constraints.Number] = g.Rect[T]

func NewRect[T constraints.Number](x, y, width, height T) Rect[T] {
	return g.NewRect(x, y, width, height)
}

// point
type Point[T constraints.Number] = g.Point2D[T]

func NewPoint[T constraints.Number](x, y T) Point[T] {
	return g.NewPoint2D(x, y)
}
