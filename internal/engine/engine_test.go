package engine

import (
	"testing"

	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/pkg"
	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/time"
)

func TestPkgs(t *testing.T) {
	settings := Settings{
		WindowWidth:  800,
		WindowHeight: 600,
		WindowTitle:  "Test Game",
	}

	game := NewGame().
		WithSettings(settings).
		WithPlugins(new(asset.AssetPlugin)).
		build()

	as, ok := core.GetResource[*asset.Server](game.game.context)
	if !ok {
		t.Fatal("expected to find asset server in game context")
	}

	if as == nil {
		t.Fatal("expected asset server to be non-nil")
	}
}

func TestDefaultPlugins(t *testing.T) {
	plugins := pkg.DefaultPlugins()
	if len(plugins) == 0 {
		t.Fatal("expected default plugins to be non-empty")
	}

	for _, pkg := range plugins {
		if pkg == nil {
			t.Fatal("expected all default plugins to be non-nil")
		}
	}

	settings := Settings{
		WindowWidth:  800,
		WindowHeight: 600,
		WindowTitle:  "Test Game",
	}

	game := NewGame().
		WithSettings(settings).
		WithPlugins(plugins...).
		build()

	as, ok := core.GetResource[*asset.Server](game.game.context)
	if !ok {
		t.Fatal("expected to find asset server in game context")
	}

	if as == nil {
		t.Fatal("expected asset server to be non-nil")
	}

	ti, ok := core.GetResource[*time.Time](game.game.context)
	if !ok {
		t.Fatal("expected to find time resource in game context")
	}

	if ti == nil {
		t.Fatal("expected time resource to be non-nil")
	}
}

// =======
// helpers
// =======
type mockResource struct{}
