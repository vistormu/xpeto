package physics

import (
	"math"

	"github.com/vistormu/xpeto/core/pkg/transform"
)

type Rect struct {
	HalfW float64
	HalfH float64
}

func (r Rect) AABB(tr *transform.Transform) AABB {
	// rotate 4 corners, then box them
	hx, hy := r.HalfW, r.HalfH
	cx, cy := +hx, +hy
	cornersX := [...]float64{+cx, +cx, -cx, -cx}
	cornersY := [...]float64{+cy, -cy, +cy, -cy}
	s, c := math.Sincos(tr.Rotation)

	worldX := make([]float64, 4)
	worldY := make([]float64, 4)
	for i := 0; i < 4; i++ {
		x, y := cornersX[i], cornersY[i]
		wx := x*c - y*s + tr.X
		wy := x*s + y*c + tr.Y
		worldX[i], worldY[i] = wx, wy
	}

	return aabbOfPoints(worldX, worldY)
}
