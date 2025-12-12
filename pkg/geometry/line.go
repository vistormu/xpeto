package geometry

import (
	c "github.com/vistormu/go-dsa/constraints"
)

type Line[T c.Number] struct {
	Point     Vector[T]
	Direction Vector[T]
}

func NewLine[T c.Number](p0, p1 Vector[T]) Line[T] {
	return Line[T]{
		Point:     p0,
		Direction: Vector[T]{X: p1.X - p0.X, Y: p1.Y - p0.Y},
	}
}

func (l Line[T]) End() Vector[T] {
	return Vector[T]{X: l.Point.X + l.Direction.X, Y: l.Point.Y + l.Direction.Y}
}
