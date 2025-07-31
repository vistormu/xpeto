package pkg

import (
	"github.com/vistormu/xpeto/internal/core"

	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/time"
)

func DefaultPlugins() []core.Plugin {
	return []core.Plugin{
		new(asset.AssetPlugin),
		new(time.TimePlugin),
	}
}
