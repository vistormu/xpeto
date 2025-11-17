package physics

import "github.com/vistormu/xpeto/pkg/transform"

// ======
// circle
// ======
type Circle struct {
	Radius float64
}

func (c Circle) AABB(tr *transform.Transform) AABB {
	return AABB{
		MinX: tr.X - c.Radius,
		MinY: tr.Y - c.Radius,
		MaxX: tr.X + c.Radius,
		MaxY: tr.Y + c.Radius,
	}
}
