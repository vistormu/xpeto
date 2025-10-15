package vector

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// ======
// circle
// ======
type Circle struct {
	Radius    float32
	Fill      color.Color
	Stroke    color.Color
	Linewidth float32
	Layer     uint16
	Order     uint16
}

type circleItem struct {
	radius    float32
	fill      color.Color
	stroke    color.Color
	linewidth float32
	x, y      float32
	layer     uint16
	order     uint16
	key       uint64
}

func (ci *circleItem) Draw(screen *ebiten.Image) {
	if ci.fill != nil && ci.radius > 0 {
		vector.DrawFilledCircle(screen, ci.x, ci.y, ci.radius, ci.fill, false)
	}
	if ci.linewidth > 0 && ci.stroke != nil && ci.radius > 0 {
		vector.StrokeCircle(screen, ci.x, ci.y, ci.radius, ci.linewidth, ci.stroke, false)
	}
}

func (ci *circleItem) SortKey() uint64 {
	return ci.key
}

// ====
// rect
// ====
type Rect struct {
	Width     float32
	Height    float32
	Fill      color.Color
	Stroke    color.Color
	Linewidth float32
	Layer     uint16
	Order     uint16
}

type rectItem struct {
	width     float32
	height    float32
	fill      color.Color
	stroke    color.Color
	linewidth float32
	x, y      float32
	layer     uint16
	order     uint16
	key       uint64
}

func (ri *rectItem) Draw(screen *ebiten.Image) {
	if ri.fill != nil && ri.width > 0 && ri.height > 0 {
		vector.DrawFilledRect(screen, ri.x, ri.y, ri.width, ri.height, ri.fill, false)
	}
	if ri.linewidth > 0 && ri.stroke != nil && ri.width > 0 && ri.height > 0 {
		vector.StrokeRect(screen, ri.x, ri.y, ri.width, ri.height, ri.linewidth, ri.stroke, false)
	}
}

func (ri *rectItem) SortKey() uint64 {
	return ri.key
}
