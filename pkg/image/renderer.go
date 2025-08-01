package image

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"

	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/transform"
)

type renderableData struct {
	position core.Vector[float32]
	scale    core.Vector[float32]
	layer    int
	image    *Image
}

type Renderer struct {
	renderables []renderableData
}

func NewRenderer() *Renderer {
	return &Renderer{
		renderables: make([]renderableData, 0),
	}
}

func (r *Renderer) Update(ctx *core.Context) {
	w := core.MustResource[*ecs.World](ctx)
	as, ok := core.GetResource[*asset.Server](ctx)
	if !ok {
		return
	}

	entities := w.Query(ecs.And(
		ecs.Has[*Renderable](),
		ecs.Has[*transform.Transform](),
	))

	// to renderable data
	renderables := make([]renderableData, 0, len(entities))
	for _, e := range entities {
		renderable, _ := ecs.GetComponent[*Renderable](w, e)
		transform, _ := ecs.GetComponent[*transform.Transform](w, e)

		img, ok := asset.GetAsset[*Image](as, renderable.Image)
		if !ok {
			continue
		}

		renderables = append(renderables, renderableData{
			position: transform.Position,
			scale:    transform.Scale,
			layer:    renderable.Layer,
			image:    img,
		})
	}

	// sort renderables by layer
	sort.Slice(renderables, func(i, j int) bool {
		return renderables[i].layer < renderables[j].layer
	})

	// update renderables
	r.renderables = renderables
}

func (r *Renderer) Draw(ctx *core.Context) {
	screen, ok := core.GetResource[*ebiten.Image](ctx)
	if !ok {
		return
	}

	for _, rend := range r.renderables {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(float64(rend.scale.X), float64(rend.scale.Y))
		opts.GeoM.Translate(float64(rend.position.X), float64(rend.position.Y))

		screen.DrawImage(rend.image.Img, opts)
	}
}
