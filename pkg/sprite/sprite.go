package sprite

import (
	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/render"
)

// ======
// sprite
// ======
type Sprite struct {
	Image    asset.Asset
	OrderKey render.OrderKey
}

func NewSprite(img asset.Asset, opts ...option) Sprite {
	s := Sprite{
		Image:    img,
		OrderKey: render.NewOrderKey(0, 0, 0),
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&s)
		}
	}

	return s
}

// =======
// options
// =======
type option func(*Sprite)

type spriteOpt struct{}

var SpriteOpt spriteOpt

func (spriteOpt) Image(a asset.Asset) option {
	return func(s *Sprite) { s.Image = a }
}

func (spriteOpt) Order(layer uint16, order uint16, tie uint32) option {
	return func(s *Sprite) { s.OrderKey = render.NewOrderKey(layer, order, tie) }
}
