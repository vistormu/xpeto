package asset

import (
	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/scheduler"
)

type AssetPkg struct{}

func (ap *AssetPkg) Resources() []any {
	return []any{
		NewServer(),
	}
}

func (ap *AssetPkg) Schedules() []*scheduler.Schedule {
	return []*scheduler.Schedule{
		{
			Name:      "AssetPkg",
			Stage:     scheduler.First,
			System:    ap.Build,
			Before:    []string{},
			After:     []string{},
			Condition: nil,
		},
	}
}

func (ap *AssetPkg) Build(ctx *core.Context) {}
