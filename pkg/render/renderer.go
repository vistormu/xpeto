package render

import (
	"image/color"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	ebitext "github.com/hajimehoshi/ebiten/v2/text"

	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"

	"github.com/vistormu/xpeto/pkg/asset"
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
}

// ========
// renderer
// ========
type Renderer struct {
	sprites []spriteData
	texts   []textData
}

func NewRenderer() *Renderer {
	return &Renderer{
		sprites: make([]spriteData, 0),
		texts:   make([]textData, 0),
	}
}

func (r *Renderer) updateSprites(ctx *core.Context) {
	w := core.MustResource[*ecs.World](ctx)
	as, ok := core.GetResource[*asset.Server](ctx)
	if !ok {
		return
	}

	entities := w.Query(ecs.And(
		ecs.Has[*image.Sprite](),
		ecs.Has[*transform.Transform](),
	))

	var sprites []spriteData
	for _, e := range entities {
		sprite, _ := ecs.GetComponent[*image.Sprite](w, e)
		transform, _ := ecs.GetComponent[*transform.Transform](w, e)

		img, ok := asset.GetAsset[*image.Image](as, sprite.Image)
		if !ok {
			continue
		}

		sprites = append(sprites, spriteData{
			position: transform.Position,
			scale:    transform.Scale,
			layer:    sprite.Layer,
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
		ecs.Has[*text.Text](),
		ecs.Has[*transform.Transform](),
	))

	var texts []textData
	for _, e := range entities {
		txt, _ := ecs.GetComponent[*text.Text](w, e)
		transform, _ := ecs.GetComponent[*transform.Transform](w, e)

		font, ok := asset.GetAsset[*text.Font](as, txt.Font)
		if !ok {
			continue
		}

		texts = append(texts, textData{
			position: transform.Position,
			layer:    txt.Layer,
			font:     font,
			content:  txt.Content,
			color:    txt.Color,
		})
	}

	sort.Slice(texts, func(i, j int) bool {
		return texts[i].layer < texts[j].layer
	})

	r.texts = texts
}

func (r *Renderer) Update(ctx *core.Context) {
	r.updateSprites(ctx)
	r.updateTexts(ctx)
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

	for _, txt := range r.texts {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(txt.position.X), float64(txt.position.Y))

		ebitext.DrawWithOptions(screen, txt.content, txt.font.Face, opts)
	}
}
