package render

import (
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	ebitext "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"

	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/graphics"
	"github.com/vistormu/xpeto/pkg/image"
	"github.com/vistormu/xpeto/pkg/text"
	"github.com/vistormu/xpeto/pkg/transform"
)

// ============
// data helpers
// ============
type spriteData struct {
	position core.Vector[float32]
	scale    core.Vector[float32]
	layer    int
	image    *image.Image
}

type textData struct {
	position core.Vector[float32]
	layer    int
	font     *text.Font
	content  string
	color    color.Color
	size     float64
}

type circleData struct {
	position  core.Vector[float32]
	layer     int
	radius    int
	fill      color.Color
	stroke    color.Color
	linewidth int
}

// ========
// renderer
// ========
type Renderer struct {
	sprites []spriteData
	texts   []textData
	circles []circleData
}

func NewRenderer() *Renderer {
	return &Renderer{
		sprites: make([]spriteData, 0),
		texts:   make([]textData, 0),
		circles: make([]circleData, 0),
	}
}

func (r *Renderer) updateSprites(ctx *core.Context) {
	w := core.MustResource[*ecs.World](ctx)
	as, ok := core.GetResource[*asset.Server](ctx)
	if !ok {
		return
	}

	entities := w.Query(ecs.And(
		ecs.Has[*Renderable](),
		ecs.Has[*image.Sprite](),
		ecs.Has[*transform.Transform](),
	))

	var sprites []spriteData
	for _, e := range entities {
		renderable, _ := ecs.GetComponent[*Renderable](w, e)
		sprite, _ := ecs.GetComponent[*image.Sprite](w, e)
		transform, _ := ecs.GetComponent[*transform.Transform](w, e)

		img, ok := asset.GetAsset[*image.Image](as, sprite.Image)
		if !ok {
			continue
		}

		if !renderable.Visible {
			continue
		}

		sprites = append(sprites, spriteData{
			position: transform.Position,
			scale:    transform.Scale,
			layer:    renderable.Layer,
			image:    img,
		})
	}

	sort.Slice(sprites, func(i, j int) bool {
		return sprites[i].layer < sprites[j].layer
	})

	r.sprites = sprites
}

func (r *Renderer) updateTexts(ctx *core.Context) {
	w := core.MustResource[*ecs.World](ctx)
	as, ok := core.GetResource[*asset.Server](ctx)
	if !ok {
		return
	}

	entities := w.Query(ecs.And(
		ecs.Has[*Renderable](),
		ecs.Has[*text.Text](),
		ecs.Has[*transform.Transform](),
	))

	var texts []textData
	for _, e := range entities {
		renderable, _ := ecs.GetComponent[*Renderable](w, e)
		txt, _ := ecs.GetComponent[*text.Text](w, e)
		transform, _ := ecs.GetComponent[*transform.Transform](w, e)

		font, ok := asset.GetAsset[*text.Font](as, txt.Font)
		if !ok {
			continue
		}

		if !renderable.Visible {
			continue
		}

		texts = append(texts, textData{
			position: transform.Position,
			layer:    renderable.Layer,
			font:     font,
			content:  txt.Content,
			color:    txt.Color,
			size:     txt.Size,
		})
	}

	sort.Slice(texts, func(i, j int) bool {
		return texts[i].layer < texts[j].layer
	})

	r.texts = texts
}

func (r *Renderer) updateCircles(ctx *core.Context) {
	w := core.MustResource[*ecs.World](ctx)

	entities := w.Query(ecs.And(
		ecs.Has[*Renderable](),
		ecs.Has[*graphics.Circle](),
		ecs.Has[*transform.Transform](),
	))

	var circles []circleData
	for _, e := range entities {
		renderable, _ := ecs.GetComponent[*Renderable](w, e)
		circle, _ := ecs.GetComponent[*graphics.Circle](w, e)
		transform, _ := ecs.GetComponent[*transform.Transform](w, e)

		if !renderable.Visible {
			continue
		}

		circles = append(circles, circleData{
			position:  transform.Position,
			layer:     renderable.Layer,
			radius:    circle.Radius,
			stroke:    circle.Stroke,
			fill:      circle.Fill,
			linewidth: circle.Linewidth,
		})
	}

	sort.Slice(circles, func(i, j int) bool {
		return circles[i].layer < circles[j].layer
	})

	r.circles = circles
}

func (r *Renderer) Update(ctx *core.Context) {
	r.updateSprites(ctx)
	r.updateTexts(ctx)
	r.updateCircles(ctx)
}

func (r *Renderer) Draw(ctx *core.Context) {
	screen, ok := core.GetResource[*ebiten.Image](ctx)
	if !ok {
		return
	}

	for _, rend := range r.sprites {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(float64(rend.scale.X), float64(rend.scale.Y))
		opts.GeoM.Translate(float64(rend.position.X), float64(rend.position.Y))

		screen.DrawImage(rend.image.Img, opts)
	}

	// TODO: sort properly the renderable stuff
	for _, txt := range r.texts {
		opts := &ebitext.DrawOptions{}
		opts.GeoM.Translate(float64(txt.position.X), float64(txt.position.Y))
		opts.ColorScale.ScaleWithColor(txt.color)
		ebitext.Draw(screen, txt.content, &ebitext.GoTextFace{
			Source: txt.font.Face,
			Size:   txt.size,
		}, opts)
	}

	for _, circle := range r.circles {
		if circle.fill != nil && circle.radius > 0 {
			vector.DrawFilledCircle(screen, float32(circle.position.X), float32(circle.position.Y), float32(circle.radius), circle.fill, false)
		}
		if circle.linewidth > 0 && circle.stroke != nil && circle.radius > 0 {
			vector.StrokeCircle(screen, float32(circle.position.X), float32(circle.position.Y), float32(circle.radius), float32(circle.linewidth), circle.stroke, false)
		}
	}
}
