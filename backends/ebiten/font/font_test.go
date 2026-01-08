package font

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/core"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/core/window"
	"github.com/vistormu/xpeto/pkg"
	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/render"
	"github.com/vistormu/xpeto/pkg/text"
	"github.com/vistormu/xpeto/pkg/transform"
)

func newTestWorld() (*ecs.World, *schedule.Scheduler) {
	w := ecs.NewWorld()
	sch := schedule.NewScheduler()

	core.CorePkgs(w, sch)
	pkg.DefaultPkgs(w, sch)
	render.Pkg[ebiten.Image](w, sch)
	Pkg(w, sch)

	return w, sch
}

// func runUpdates(w *ecs.World, sch *schedule.Scheduler, n int) {
// 	for i := 0; i < n; i++ {
// 		schedule.RunUpdate(w, sch)
// 	}
// }

func TestExtractText_DefaultFontNotLoaded_SkipsSilently(t *testing.T) {
	w := ecs.NewWorld()

	e := ecs.AddEntity(w)
	ecs.AddComponent(w, e, text.NewText("hello")) // Font == 0
	ecs.AddComponent(w, e, transform.Transform{X: 10, Y: 20})

	got := extractText(w)
	if len(got) != 0 {
		t.Fatalf(
			"expected no renderables when default font is not loaded, got %d",
			len(got),
		)
	}
}

func TestExtractText_DefaultFontLoaded_Renders(t *testing.T) {
	w, sch := newTestWorld()

	loaded := false
	for range 240 {
		schedule.RunUpdate(w, sch)
		if _, ok := text.GetDefaultFont(w); ok {
			loaded = true
			break
		}
	}
	if !loaded {
		t.Fatalf("default font did not load; check assets.DefaultFS and loader registration")
	}

	e := ecs.AddEntity(w)
	ecs.AddComponent(w, e, text.NewText("hello")) // Font == 0 → default
	ecs.AddComponent(w, e, transform.Transform{X: 1, Y: 2})

	got := extractText(w)
	if len(got) != 1 {
		t.Fatalf("expected 1 renderable, got %d", len(got))
	}
	if got[0].content != "hello" {
		t.Fatalf("expected content 'hello', got %q", got[0].content)
	}
	if got[0].face == nil {
		t.Fatalf("expected face to be non nil")
	}
}

func TestExtractText_MissingExplicitFont_Skips(t *testing.T) {
	w := ecs.NewWorld()
	ecs.AddResource(w, window.Scaling{SnapPixels: false})

	e := ecs.AddEntity(w)
	txt := text.NewText(
		"hello",
		text.TextOpt.Font(asset.Asset(123456)), // explicit, missing
	)
	ecs.AddComponent(w, e, txt)
	ecs.AddComponent(w, e, transform.Transform{X: 0, Y: 0})

	got := extractText(w)
	if len(got) != 0 {
		t.Fatalf(
			"expected no renderables for missing explicit font, got %d",
			len(got),
		)
	}
}
