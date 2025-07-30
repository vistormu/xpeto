package pkg

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/pkg/asset"
)

type DefaultPlugins struct{}

func (dp *DefaultPlugins) Resources() []any {
	return []any{}
}

func (dp *DefaultPlugins) Systems() []core.System {
	return []core.System{}
}

func (dp *DefaultPlugins) Build(ctx *core.Context) {}
