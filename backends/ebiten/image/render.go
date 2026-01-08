package image

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/backends/ebiten/shared"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/window"
	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/render"
	"github.com/vistormu/xpeto/pkg/sprite"
	"github.com/vistormu/xpeto/pkg/transform"
)

type renderable struct {
	img      *image
	x, y     float64
	anchor   render.AnchorType
	snap     bool
	orderKey render.OrderKey
}

type renderableBuffer struct {
	renderables []renderable
}

func newRenderableBuffer() renderableBuffer {
	return renderableBuffer{
		renderables: make([]renderable, 0),
	}
}

func extractSprite(w *ecs.World) []renderable {
	buf := ecs.EnsureResource(w, newRenderableBuffer)
	buf.renderables = buf.renderables[:0]

	sc, ok := ecs.GetResource[window.Scaling](w)
	snap := ok && sc.SnapPixels

	q := ecs.NewQuery2[sprite.Sprite, transform.Transform](w)

	for _, b := range q.Iter() {
		s, tr := b.Components()
		e := b.Entity()

		img, ok := asset.GetAsset[image](w, s.Image)
		if !ok || img == nil || img.Image == nil {
			continue
		}

		anchor := render.AnchorCenter
		if an, ok := ecs.GetComponent[render.Anchor](w, e); ok && an != nil {
			anchor = an.Type
		}

		buf.renderables = append(buf.renderables, renderable{
			img:      img,
			x:        tr.X,
			y:        tr.Y,
			anchor:   anchor,
			snap:     snap,
			orderKey: s.OrderKey,
		})
	}

	return buf.renderables
}

func sortSprite(s renderable) uint64 {
	return uint64(s.orderKey)
}

func drawSprite(screen *ebiten.Image, s renderable) {
	w := float64(s.img.Bounds().Dx())
	h := float64(s.img.Bounds().Dy())

	dx, dy := shared.Offset(w, h, s.anchor)

	x := s.x + dx
	y := s.y + dy

	if s.snap {
		x = math.Round(x)
		y = math.Round(y)
	}

	var op ebiten.DrawImageOptions
	op.GeoM.Translate(x, y)
	screen.DrawImage(s.img.Image, &op)
}
