package text

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/pkg/asset"
)

type DefaultFont struct {
	Font asset.Asset `path:"default/font.ttf"`
}

func GetDefaultFont(w *ecs.World) (asset.Asset, bool) {
	df, ok := ecs.GetResource[DefaultFont](w)
	if !ok {
		return asset.Asset(0), false
	}

	ok = asset.IsAssetLoaded(w, df.Font)
	if !ok {
		return asset.Asset(0), false
	}

	return df.Font, true
}
