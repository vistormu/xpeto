package geometry

import (
	c "github.com/vistormu/go-dsa/constraints"
)

type Path[T c.Number] struct {
	Points []Vector[T]
	Closed bool
}

func NewPath[T c.Number]() Path[T] {
	return Path[T]{
		Points: make([]Vector[T], 0),
		Closed: false,
	}
}

func (p *Path[T]) AddPoint(x, y T) {
	p.Points = append(p.Points, Vector[T]{x, y})
}
