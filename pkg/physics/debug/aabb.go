package debug

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/pkg/physics"
	"github.com/vistormu/xpeto/pkg/transform"
)

type aabb struct {
	x, y, w, h float32
	border     float32
	color      color.Color
}

func extractAabb(w *ecs.World) []aabb {
	s, _ := ecs.GetResource[Settings](w)

	items := make([]aabb, 0)

	if !s.Enabled || !s.DrawAABBs {
		return items
	}

	q := ecs.NewQuery2[physics.AABB, transform.Transform](w)
	for _, b := range q.Iter() {
		a, tr := b.Components()

		items = append(items, aabb{
			x:      float32(tr.X - (a.MaxX-a.MinX)/2),
			y:      float32(tr.Y - (a.MaxY-a.MinY)/2),
			w:      float32(a.MaxX - a.MinX),
			h:      float32(a.MaxY - a.MinY),
			border: s.LineWidthPx,
			color:  s.AABBStroke,
		})
	}

	return items

}

func sortAabb(a aabb) uint64 {
	return 0
}

func drawAabb(screen *ebiten.Image, a aabb) {
	if a.w > 0 && a.h > 0 && a.border > 0 {
		vector.StrokeRect(screen, a.x, a.y, a.w, a.h, a.border, a.color, false)
	}
}
