package pkg

import (
	"github.com/vistormu/xpeto/internal/game"

	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/font"
	"github.com/vistormu/xpeto/pkg/image"
	"github.com/vistormu/xpeto/pkg/input"
	"github.com/vistormu/xpeto/pkg/render"
	"github.com/vistormu/xpeto/pkg/sprite"
	"github.com/vistormu/xpeto/pkg/text"
	"github.com/vistormu/xpeto/pkg/time"
	"github.com/vistormu/xpeto/pkg/vector"
)

func DefaultPlugins() []game.Plugin {
	return []game.Plugin{
		// core: no dependencies
		asset.AssetPlugin,
		time.TimePlugin,
		input.InputPlugin,

		// semi-core: depends only on core
		image.ImagePlugin,
		render.RenderPlugin,
		font.FontPlugin,

		// features: depends on core and semi-core
		sprite.SpritePlugin,
		text.TextPlugin,
		vector.VectorPlugin,
	}
}
