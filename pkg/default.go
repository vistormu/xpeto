package pkg

import (
	"github.com/vistormu/xpeto/internal/game"

	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/image"
	"github.com/vistormu/xpeto/pkg/input"
	"github.com/vistormu/xpeto/pkg/render"
	"github.com/vistormu/xpeto/pkg/text"
	"github.com/vistormu/xpeto/pkg/time"
)

func DefaultPlugins() []game.Plugin {
	return []game.Plugin{
		asset.AssetPlugin,
		time.TimePlugin,
		image.ImagePlugin,
		text.TextPlugin,
		render.RenderPlugin,
		input.InputPlugin,
	}
}
