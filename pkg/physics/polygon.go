package physics

import (
	"math"
)

type Convex struct {
	XY []float64
}

func (p Convex) AABB(px, py, rot float64) AABB {
	n := len(p.XY) / 2
	s, c := math.Sincos(rot)
	xs := make([]float64, n)
	ys := make([]float64, n)
	for i := 0; i < n; i++ {
		x := p.XY[2*i]
		y := p.XY[2*i+1]
		wx := x*c - y*s + px
		wy := x*s + y*c + py
		xs[i], ys[i] = wx, wy
	}
	return aabbOfPoints(xs, ys)
}
