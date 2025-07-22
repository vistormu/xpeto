package render

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"
	"github.com/vistormu/xpeto/pkg/transform"
)

type renderableData struct {
	position core.Vector[float32]
	scale    core.Vector[float32]
	renderer *Renderer
	layer    int
}

type System struct {
	renderables []renderableData
}

func NewSystem() *System {
	return &System{
		renderables: make([]renderableData, 0),
	}
}

func (r *System) Update(ctx *ecs.Context) {
	em, _ := ecs.GetResource[*ecs.Manager](ctx)
	im, _ := ecs.GetResource[*Manager](ctx)

	entities := em.Query(ecs.And(
		ecs.Has[*Renderable](),
		ecs.Has[*transform.Transform](),
	))

	// to renderable data
	renderables := make([]renderableData, 0, len(entities))
	for _, e := range entities {
		renderable, _ := ecs.GetComponent[*Renderable](em, e)
		transform, _ := ecs.GetComponent[*transform.Transform](em, e)

		renderables = append(renderables, renderableData{
			position: transform.Position,
			scale:    transform.Scale,
			renderer: im.Renderer(renderable.Image),
			layer:    renderable.Layer,
		})
	}

	// sort renderables by layer
	sort.Slice(renderables, func(i, j int) bool {
		return renderables[i].layer < renderables[j].layer
	})

	// update renderables
	r.renderables = renderables
}

func (r *System) Draw(screen *Renderer) {
	for _, rend := range r.renderables {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(float64(rend.scale.X), float64(rend.scale.Y))
		opts.GeoM.Translate(float64(rend.position.X), float64(rend.position.Y))

		screen.DrawImage(rend.renderer, opts)
	}
}
