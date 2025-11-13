package pkg

import (
	"github.com/vistormu/xpeto/core/pkg/event"
	"github.com/vistormu/xpeto/core/pkg/log"
	"github.com/vistormu/xpeto/core/pkg/time"
	"github.com/vistormu/xpeto/core/pkg/window"
)

func CorePkgs() []Pkg {
	return []Pkg{
		event.Pkg,
		time.Pkg,
		window.Pkg,
		log.Pkg,
	}
}
