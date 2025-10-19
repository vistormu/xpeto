package physics

import (
	"math"
)

type Segment struct {
	AX, AY float64
	BX, BY float64
}

func (s Segment) AABB(px, py, rot float64) AABB {
	sn, cs := math.Sincos(rot)
	ax := s.AX*cs - s.AY*sn + px
	ay := s.AX*sn + s.AY*cs + py
	bx := s.BX*cs - s.BY*sn + px
	by := s.BX*sn + s.BY*cs + py

	return aabbOfPoints([]float64{ax, bx}, []float64{ay, by})
}
