package debug

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/pkg/transform"
	"github.com/vistormu/xpeto/pkg/physics"
	"github.com/vistormu/xpeto/pkg/render"
)

// ====
// item
// ====
type aabbItem struct {
	x, y, w, h float32
	border     float32
	color      color.Color
	key        uint64
}

func (ri *aabbItem) Draw(screen *ebiten.Image) {
	if ri.w > 0 && ri.h > 0 && ri.border > 0 {
		vector.StrokeRect(screen, ri.x, ri.y, ri.w, ri.h, ri.border, ri.color, false)
	}
}
func (ri *aabbItem) SortKey() uint64 { return ri.key }

// =========
// extractor
// =========
func extractAabb(w *ecs.World) []render.Renderable {
	s, _ := ecs.GetResource[Settings](w)

	items := make([]render.Renderable, 0)

	if !s.Enabled || !s.DrawAABBs {
		return items
	}

	key := packKey(s.Layer, s.Order)

	q := ecs.NewQuery2[physics.AABB, transform.Transform](w)
	for _, b := range q.Iter() {
		a := b.A()
		tr := b.B()

		items = append(items, &aabbItem{
			x:      float32(tr.X - (a.MaxX-a.MinX)/2),
			y:      float32(tr.Y - (a.MaxY-a.MinY)/2),
			w:      float32(a.MaxX - a.MinX),
			h:      float32(a.MaxY - a.MinY),
			border: s.LineWidthPx,
			color:  s.AABBStroke,
			key:    key,
		})
	}

	return items
}
