package physics

import (
	"math"

	dsamath "github.com/vistormu/go-dsa/math"
	"github.com/vistormu/go-dsa/set"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/pkg/transform"
)

type AABB struct {
	MinX, MinY float64
	MaxX, MaxY float64
}

func aabbOfPoints(xs, ys []float64) AABB {
	minX, minY := math.Inf(1), math.Inf(1)
	maxX, maxY := math.Inf(-1), math.Inf(-1)
	for i := range xs {
		x, y := xs[i], ys[i]
		if x < minX {
			minX = x
		}
		if y < minY {
			minY = y
		}
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
	}

	return AABB{MinX: minX, MinY: minY, MaxX: maxX, MaxY: maxY}
}

func aabbToCellSpan(s *Space, aabb AABB) (i0, j0, i1, j1 int) {
	if s.CellWidth <= 0 || s.CellHeight <= 0 || s.cols == 0 || s.rows == 0 {
		return 0, 0, -1, -1
	}

	// reject outside
	if aabb.MaxX <= 0 || aabb.MaxY <= 0 || aabb.MinX >= s.Width || aabb.MinY >= s.Height {
		return 0, 0, -1, -1
	}

	// clamp to world bounds first (optional but helps precision)
	if aabb.MinX < 0 {
		aabb.MinX = 0
	}
	if aabb.MinY < 0 {
		aabb.MinY = 0
	}
	if aabb.MaxX > s.Width {
		aabb.MaxX = s.Width
	}
	if aabb.MaxY > s.Height {
		aabb.MaxY = s.Height
	}

	floorDiv := func(v, d float64) int { return int(math.Floor(v / d)) }

	const eps = 1e-7

	i0 = dsamath.Clip(floorDiv(aabb.MinX, s.CellWidth), 0, s.cols-1)
	j0 = dsamath.Clip(floorDiv(aabb.MinY, s.CellHeight), 0, s.rows-1)
	i1 = dsamath.Clip(floorDiv(aabb.MaxX-eps, s.CellWidth), 0, s.cols-1)
	j1 = dsamath.Clip(floorDiv(aabb.MaxY-eps, s.CellHeight), 0, s.rows-1)

	if i0 > i1 || j0 > j1 {
		return 0, 0, -1, -1
	}

	return
}

func fillGrid(w *ecs.World) {
	s, _ := ecs.GetResource[Space](w)

	s.Clear()

	q := ecs.NewQuery2[Collider, transform.Transform](w)

	for _, b := range q.Iter() {
		col := b.A()
		tr := b.B()

		i0, j0, i1, j1 := aabbToCellSpan(s, col.Shape.AABB(tr))
		if i1 < i0 || j1 < j0 {
			continue
		}

		for j := j0; j <= j1; j++ {
			for i := i0; i <= i1; i++ {
				s.AddEntity(b.Entity(), i, j)
			}
		}
	}
}

func getCandidates(w *ecs.World) {
	s, _ := ecs.GetResource[Space](w)

	seen := set.NewHashSet[uint64]()

	neighbors := [][2]int{
		{0, 0}, {1, 0}, {0, 1}, {1, 1},
		{-1, 0}, {0, -1}, {-1, -1}, {1, -1}, {-1, 1},
	}

	pass := func(a, b *Collider) bool {
		return (a.Mask&b.Layer) != 0 || (b.Mask&a.Layer) != 0
	}

	for j := 0; j < s.rows; j++ {
		for i := 0; i < s.cols; i++ {
			base, ok := s.GetCell(i, j)
			if !ok || s.IsEmpty(i, j) {
				continue
			}

			for _, off := range neighbors {
				nc, ok := s.GetCell(i+off[0], j+off[1])
				if !ok || s.IsEmpty(i+off[0], j+off[1]) {
					continue
				}

				for _, ea := range base.Entities {
					for _, eb := range nc.Entities {
						if ea == eb {
							continue
						}
						// order pair
						a, b := ea, eb
						if b < a {
							a, b = b, a
						}

						key := (uint64(a) << 32) | uint64(b)
						if seen.Contains(key) {
							continue
						}

						// layer/mask prefilter
						ca, okA := ecs.GetComponent[Collider](w, a)
						cb, okB := ecs.GetComponent[Collider](w, b)
						if !okA || !okB || !pass(ca, cb) {
							continue
						}

						seen.Add(key)
						s.candidates = append(s.candidates, pair{A: a, B: b})
					}
				}
			}
		}
	}
}
