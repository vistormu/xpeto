package image

import (
	"github.com/vistormu/xpeto/pkg/asset"
)

type Sprite struct {
	Image asset.Handle
	Layer int
	FlipX bool
	FlipY bool
}
