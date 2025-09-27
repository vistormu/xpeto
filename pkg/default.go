package pkg

import (
	"github.com/vistormu/xpeto/internal/core"

	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/image"
	"github.com/vistormu/xpeto/pkg/render"
	"github.com/vistormu/xpeto/pkg/text"
	"github.com/vistormu/xpeto/pkg/time"
)

func DefaultPlugins() []core.Plugin {
	return []core.Plugin{
		asset.AssetPlugin,
		time.TimePlugin,
		image.ImagePlugin,
		text.TextPlugin,
		render.RenderPlugin,
	}
}
