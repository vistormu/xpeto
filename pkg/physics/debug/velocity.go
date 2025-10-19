package debug

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/pkg/transform"

	"github.com/vistormu/xpeto/pkg/physics"
	"github.com/vistormu/xpeto/pkg/render"
)

func packKey(layer, order uint16) uint64 { return (uint64(layer) << 16) | uint64(order) }

// ====
// item
// ====
type velocityItem struct {
	x1, y1, x2, y2 float32
	w              float32
	color          color.Color
	key            uint64
}

func (li *velocityItem) Draw(screen *ebiten.Image) {
	// body
	vector.StrokeLine(screen, li.x1, li.y1, li.x2, li.y2, li.w, li.color, false)

	// head
	dx, dy := li.x2-li.x1, li.y2-li.y1
	l := float32(math.Hypot(float64(dx), float64(dy)))

	if l <= 0 {
		return
	}

	ax := dx / l
	ay := dy / l
	px := -ay
	py := ax
	head := float32(6.0)
	vector.StrokeLine(screen,
		li.x2,
		li.y2,
		li.x2-ax*head+px*head*0.4,
		li.y2-ay*head+py*head*0.4,
		li.w,
		li.color,
		false,
	)
	vector.StrokeLine(screen,
		li.x2,
		li.y2,
		li.x2-ax*head-px*head*0.4,
		li.y2-ay*head-py*head*0.4,
		li.w,
		li.color,
		false,
	)
}

func (li *velocityItem) SortKey() uint64 { return li.key }

// =========
// extractor
// =========
func extractVelocities(w *ecs.World) []render.Renderable {
	s, _ := ecs.GetResource[Settings](w)

	items := make([]render.Renderable, 0)

	if !s.Enabled || !s.DrawVelocities {
		return items
	}

	key := packKey(s.Layer, s.Order)

	q := ecs.NewQuery2[physics.Velocity, transform.Transform](w)
	for _, b := range q.Iter() {
		v := b.A()
		tr := b.B()
		x1 := float32(tr.X)
		y1 := float32(tr.Y)
		x2 := x1 + float32(v.X)*s.VelocityScale
		y2 := y1 + float32(v.Y)*s.VelocityScale

		items = append(items, &velocityItem{
			x1: x1, y1: y1, x2: x2, y2: y2,
			w: s.LineWidthPx, color: s.VelocityColor, key: key,
		})
	}

	return items
}
