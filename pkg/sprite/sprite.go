package sprite

import (
	"github.com/vistormu/xpeto/pkg/asset"
)

type Sprite struct {
	Image asset.Asset
	FlipX bool
	FlipY bool
	Layer uint16
	Order uint16
}
