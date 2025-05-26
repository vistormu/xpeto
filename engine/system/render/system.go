package render

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/ecs"
	"github.com/vistormu/xpeto/engine/component"
	"github.com/vistormu/xpeto/image"
	g "github.com/vistormu/xpeto/internal/geometry"
)

type renderableData struct {
	position g.Vector[float32]
	scale    g.Vector[float32]
	image    *image.Image
	layer    int
}

type System struct {
	renderables []renderableData
}

func NewSystem() ecs.System {
	return &System{
		renderables: make([]renderableData, 0),
	}
}

func (s *System) OnLoad(ctx *ecs.Context)           {}
func (s *System) OnUnload(*ecs.Context)             {}
func (s *System) FixedUpdate(*ecs.Context, float32) {}

func (r *System) Update(ctx *ecs.Context, dt float32) {
	em, _ := ecs.GetResource[*ecs.EntityManager](ctx)
	im, _ := ecs.GetResource[*image.Manager](ctx)

	entities := em.Query(ecs.And(
		ecs.Has[*component.Renderable](),
		ecs.Has[*component.Transform](),
	))

	// to renderable data
	renderables := make([]renderableData, 0, len(entities))
	for _, e := range entities {
		renderable, _ := ecs.GetComponent[*component.Renderable](em, e)
		transform, _ := ecs.GetComponent[*component.Transform](em, e)

		renderables = append(renderables, renderableData{
			position: transform.Position,
			scale:    transform.Scale,
			image:    im.Image(renderable.Sprite),
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

func (r *System) Draw(screen *image.Image) {
	for _, rend := range r.renderables {
		screen.DrawImage(rend.image, &ebiten.DrawImageOptions{})
	}
}
