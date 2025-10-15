package pkg

import (
	"github.com/vistormu/xpeto/core/pkg"

	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/font"
	"github.com/vistormu/xpeto/pkg/image"
	"github.com/vistormu/xpeto/pkg/input"
	"github.com/vistormu/xpeto/pkg/render"
	"github.com/vistormu/xpeto/pkg/sprite"
	"github.com/vistormu/xpeto/pkg/text"
	"github.com/vistormu/xpeto/pkg/vector"
)

func DefaultPkgs() []pkg.Pkg {
	return []pkg.Pkg{
		// core: no dependencies
		asset.Pkg,
		input.Pkg,

		// semi-core: depends only on core
		image.Pkg,
		render.Pkg,
		font.Pkg,

		// features: depends on core and semi-core
		sprite.Pkg,
		text.Pkg,
		vector.Pkg,
	}
}
