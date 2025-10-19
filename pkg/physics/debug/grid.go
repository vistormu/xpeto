package debug

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/xpeto/core/ecs"

	"github.com/vistormu/xpeto/pkg/physics"
	"github.com/vistormu/xpeto/pkg/render"
)

// ====
// item
// ====
type gridItem struct {
	x, y   float32
	w, h   float32
	color  color.Color
	border float32
	key    uint64
}

func (li *gridItem) Draw(screen *ebiten.Image) {
	if li.w > 0 && li.h > 0 && li.border > 0 {
		vector.StrokeRect(screen, li.x, li.y, li.w, li.h, li.border, li.color, false)
	}
}

func (li *gridItem) SortKey() uint64 { return li.key }

// =========
// extractor
// =========
func extractGrid(w *ecs.World) []render.Renderable {
	s, _ := ecs.GetResource[Settings](w)
	sp, _ := ecs.GetResource[physics.Space](w)
	cells := sp.Cells

	items := make([]render.Renderable, 0)

	if !s.Enabled || !s.DrawOccupiedGrid {
		return items
	}

	key := packKey(s.Layer, s.Order)

	for _, c := range cells {
		x := float32(c.I) * float32(sp.CellWidth)
		y := float32(c.J) * float32(sp.CellHeight)

		if sp.IsEmpty(c.I, c.J) {
			continue
		}

		items = append(items, &gridItem{
			x:      x,
			y:      y,
			w:      float32(sp.CellWidth),
			h:      float32(sp.CellHeight),
			border: s.LineWidthPx,
			color:  s.GridStroke,
			key:    key,
		})
	}

	return items
}
