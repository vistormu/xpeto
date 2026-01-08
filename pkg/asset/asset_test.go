package asset

import (
	"errors"
	"fmt"
	"testing"
	"testing/fstest"
	"time"

	"github.com/vistormu/xpeto/core"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func initWorld(t *testing.T) (*ecs.World, *schedule.Scheduler) {
	t.Helper()

	w := ecs.NewWorld()
	sch := schedule.NewScheduler()

	core.CorePkgs(w, sch)
	Pkg(w, sch)

	// some plugins register things on startup
	schedule.RunStartup(w, sch)

	return w, sch
}

func tickAll(w *ecs.World, sch *schedule.Scheduler) {
	// drive scheduler like runtime
	schedule.RunUpdate(w, sch)

	// keep tests stable even if schedule wiring changes
	readRequests(w)
	loadResults(w)
}

func waitUntil(t *testing.T, d time.Duration, step func(), pred func() bool) {
	t.Helper()

	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		step()
		if pred() {
			return
		}
		time.Sleep(1 * time.Millisecond)
	}
	t.Fatalf("timeout waiting for condition")
}

func TestSplitPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in       string
		wantBase string
		wantRel  string
		wantExt  string
	}{
		{in: "assets/images/logo.png", wantBase: "assets", wantRel: "images/logo.png", wantExt: ".png"},
		{in: "audio/sfx/click.WAV", wantBase: "audio", wantRel: "sfx/click.WAV", wantExt: ".wav"},
		{in: "file.txt", wantBase: "", wantRel: "file.txt", wantExt: ".txt"},
		{in: "noext", wantBase: "", wantRel: "noext", wantExt: ""},
		// current splitPath behaviour. if you later normalize this, update these cases.
		{in: "", wantBase: "", wantRel: ".", wantExt: "."},
		{in: "/abs/path.txt", wantBase: "", wantRel: "abs/path.txt", wantExt: ".txt"},
	}

	for _, tt := range tests {
		base, rel, ext := splitPath(tt.in)
		if base != tt.wantBase || rel != tt.wantRel || ext != tt.wantExt {
			t.Fatalf("splitPath(%q) = (%q,%q,%q), want (%q,%q,%q)",
				tt.in, base, rel, ext, tt.wantBase, tt.wantRel, tt.wantExt)
		}
	}
}

func TestAddStaticFS_BaseBehaviour(t *testing.T) {
	t.Parallel()

	w, _ := initWorld(t)
	s, _ := ecs.GetResource[server](w)

	// Pkg registers "default" FS
	if _, ok := s.staticFS["default"]; !ok {
		t.Fatal(`expected default filesystem registered under base "default"`)
	}
	prev := len(s.staticFS)

	// empty base is rejected (no registration, map unchanged)
	AddStaticFS(w, "", fstest.MapFS{})
	if len(s.staticFS) != prev {
		t.Fatalf("expected no new filesystem for empty base, got %d->%d", prev, len(s.staticFS))
	}

	// base with path separators is rejected (map unchanged)
	AddStaticFS(w, "bad/name", fstest.MapFS{})
	if len(s.staticFS) != prev {
		t.Fatalf("expected no new filesystem for base with slash, got %d->%d", prev, len(s.staticFS))
	}

	// valid base adds one entry
	AddStaticFS(w, "assets", fstest.MapFS{})
	if _, ok := s.staticFS["assets"]; !ok {
		t.Fatalf("filesystem with base %q not found", "assets")
	}
}

type testBundle struct {
	Text Asset `path:"assets/text/hello.txt"`
}

type testBundle2 struct {
	A Asset `path:"assets/text/a.txt"`
	B Asset `path:"assets/text/b.txt"`
}

type testBundleMissingTag struct {
	Text Asset
}

type testBundleBadBase struct {
	Text Asset `path:"bad/name/text/hello.txt"`
}

type testBundleUnknownBase struct {
	Text Asset `path:"missing/text/hello.txt"`
}

type testBundleMissingLoader struct {
	Text Asset `path:"assets/text/hello.txt"`
}

type testType struct {
	Value string
}

func testLoader(data []byte, _ string) (*testType, error) {
	return &testType{Value: string(data)}, nil
}

func errLoader(err error) LoaderFn[testType] {
	return func(_ []byte, _ string) (*testType, error) {
		return nil, err
	}
}

func slowLoader(delay time.Duration) LoaderFn[testType] {
	return func(data []byte, _ string) (*testType, error) {
		time.Sleep(delay)
		return &testType{Value: string(data)}, nil
	}
}

func TestAssetLoadEndToEnd(t *testing.T) {
	w, sch := initWorld(t)

	AddStaticFS(w, "assets", fstest.MapFS{
		"text/hello.txt": &fstest.MapFile{Data: []byte("hello, world"), Mode: 0o444},
	})
	AddLoaderFn(w, testLoader, ".txt")

	AddAsset[testBundle](w)

	bundle, ok := ecs.GetResource[testBundle](w)
	if !ok {
		t.Fatal("testBundle resource not found after AddAsset")
	}
	if bundle.Text == Asset(0) {
		t.Fatal("asset handle for Text is zero, expected non-zero")
	}

	// path introspection should work for non-zero handles
	if p, ok := GetAssetPath(w, bundle.Text); !ok || p != "assets/text/hello.txt" {
		t.Fatalf("GetAssetPath = (%q,%v), want (%q,true)", p, ok, "assets/text/hello.txt")
	}

	waitUntil(t, 2*time.Second,
		func() { tickAll(w, sch) },
		func() bool { return IsAssetLoaded(w, bundle.Text) },
	)

	got, ok := GetAsset[testType](w, bundle.Text)
	if !ok || got == nil {
		st, _ := GetAssetState(w, bundle.Text)
		t.Fatalf("asset state is %v but GetAsset returned (nil,false)", st)
	}
	if got.Value != "hello, world" {
		t.Fatalf("loaded asset value = %q, want %q", got.Value, "hello, world")
	}

	// RemoveAsset after load must succeed
	if !RemoveAsset[testType](w, bundle.Text) {
		t.Fatal("RemoveAsset returned false, expected true")
	}
	if v, ok := GetAsset[testType](w, bundle.Text); ok || v != nil {
		t.Fatal("expected asset to be removed, but GetAsset still returned a value")
	}

	// After removal, state queries should not report "loaded"
	if st, ok := GetAssetState(w, bundle.Text); ok && st == AssetLoaded {
		t.Fatal("expected removed asset to not remain in loaded state")
	}
}

func TestAddAsset_MissingPathTagKeepsZeroHandleButRegistersBundle(t *testing.T) {
	t.Parallel()

	w, _ := initWorld(t)

	AddAsset[testBundleMissingTag](w)

	b, ok := ecs.GetResource[testBundleMissingTag](w)
	if !ok {
		t.Fatal("expected bundle resource to exist even if it cannot be populated")
	}
	if b.Text != Asset(0) {
		t.Fatal("expected Asset handle to remain zero when path tag is missing")
	}
}

func TestAddAsset_BadBaseKeepsZeroHandleButRegistersBundle(t *testing.T) {
	t.Parallel()

	w, _ := initWorld(t)

	AddAsset[testBundleBadBase](w)

	b, ok := ecs.GetResource[testBundleBadBase](w)
	if !ok {
		t.Fatal("expected bundle resource to exist even if base is invalid")
	}
	if b.Text != Asset(0) {
		t.Fatal("expected Asset handle to remain zero for invalid base")
	}
}

func TestAddAsset_UnknownBaseKeepsZeroHandle(t *testing.T) {
	t.Parallel()

	w, _ := initWorld(t)

	AddLoaderFn(w, testLoader, ".txt") // no FS for "missing"

	AddAsset[testBundleUnknownBase](w)
	b, ok := ecs.GetResource[testBundleUnknownBase](w)
	if !ok {
		t.Fatal("expected bundle resource to exist")
	}
	if b.Text != Asset(0) {
		t.Fatal("expected Asset handle to remain zero when base FS is missing")
	}
}

func TestAddAsset_MissingLoaderKeepsZeroHandle(t *testing.T) {
	t.Parallel()

	w, _ := initWorld(t)

	AddStaticFS(w, "assets", fstest.MapFS{
		"text/hello.txt": &fstest.MapFile{Data: []byte("hi")},
	})
	// no loader registered for ".txt"

	AddAsset[testBundleMissingLoader](w)
	b, ok := ecs.GetResource[testBundleMissingLoader](w)
	if !ok {
		t.Fatal("expected bundle resource to exist")
	}
	if b.Text != Asset(0) {
		t.Fatal("expected Asset handle to remain zero when loader is missing")
	}
}

func TestAsset_LoaderErrorFailsAndExposesError(t *testing.T) {
	w, sch := initWorld(t)

	AddStaticFS(w, "assets", fstest.MapFS{
		"text/hello.txt": &fstest.MapFile{Data: []byte("hi")},
	})
	AddLoaderFn(w, errLoader(errors.New("decode failed")), ".txt")

	AddAsset[testBundle](w)
	bundle, ok := ecs.GetResource[testBundle](w)
	if !ok {
		t.Fatal("testBundle resource not found after AddAsset")
	}
	if bundle.Text == Asset(0) {
		t.Fatal("asset handle for Text is zero, expected non-zero")
	}

	waitUntil(t, 2*time.Second,
		func() { tickAll(w, sch) },
		func() bool {
			st, ok := GetAssetState(w, bundle.Text)
			return ok && st == AssetFailed
		},
	)

	if err, ok := GetAssetError(w, bundle.Text); !ok || err == nil {
		t.Fatal("expected GetAssetError to return a non-nil error after failure")
	}
}

func TestAsset_ReadFileErrorFailsAndExposesError(t *testing.T) {
	w, sch := initWorld(t)

	// FS exists, file does not
	AddStaticFS(w, "assets", fstest.MapFS{})
	AddLoaderFn(w, testLoader, ".txt")

	AddAsset[testBundle](w)
	bundle, ok := ecs.GetResource[testBundle](w)
	if !ok {
		t.Fatal("testBundle resource not found after AddAsset")
	}
	if bundle.Text == Asset(0) {
		t.Fatal("asset handle for Text is zero, expected non-zero")
	}

	waitUntil(t, 2*time.Second,
		func() { tickAll(w, sch) },
		func() bool {
			st, ok := GetAssetState(w, bundle.Text)
			return ok && st == AssetFailed
		},
	)

	err, ok := GetAssetError(w, bundle.Text)
	if !ok || err == nil {
		t.Fatal("expected read failure to be visible via GetAssetError")
	}
}

func TestAsset_RemoveBeforeLoad_DoesNotResurrect(t *testing.T) {
	w, sch := initWorld(t)

	AddStaticFS(w, "assets", fstest.MapFS{
		"text/hello.txt": &fstest.MapFile{Data: []byte("hi")},
	})
	AddLoaderFn(w, slowLoader(20*time.Millisecond), ".txt")

	AddAsset[testBundle](w)
	bundle, ok := ecs.GetResource[testBundle](w)
	if !ok {
		t.Fatal("testBundle resource not found after AddAsset")
	}
	if bundle.Text == Asset(0) {
		t.Fatal("asset handle for Text is zero, expected non-zero")
	}

	// remove quickly, before load completes
	if !RemoveAsset[testType](w, bundle.Text) {
		t.Fatal("expected RemoveAsset to succeed even if asset is not loaded yet")
	}

	// drive the system for a bit; if stale results are committed, the asset can reappear
	for i := 0; i < 100; i++ {
		tickAll(w, sch)
		time.Sleep(1 * time.Millisecond)
	}

	if v, ok := GetAsset[testType](w, bundle.Text); ok || v != nil {
		t.Fatal("asset resurrected after removal; stale result was committed")
	}
}

func TestConditions_CommonUsage(t *testing.T) {
	w, sch := initWorld(t)

	AddStaticFS(w, "assets", fstest.MapFS{
		"text/hello.txt": &fstest.MapFile{Data: []byte("hi")},
	})
	AddLoaderFn(w, testLoader, ".txt")

	AddAsset[testBundle](w)
	bundle, ok := ecs.GetResource[testBundle](w)
	if !ok {
		t.Fatal("testBundle resource not found after AddAsset")
	}

	if WhenAssetLoaded(bundle.Text)(w) {
		t.Fatal("WhenAssetLoaded should be false before load")
	}
	if WhenBundleLoaded[testBundle]()(w) {
		t.Fatal("WhenBundleLoaded should be false before load")
	}

	waitUntil(t, 2*time.Second,
		func() { tickAll(w, sch) },
		func() bool { return WhenBundleLoaded[testBundle]()(w) },
	)

	if !WhenAllAssetsLoaded(bundle.Text)(w) {
		t.Fatal("WhenAllAssetsLoaded should be true after load")
	}
	if WhenAnyAssetFailed(bundle.Text)(w) {
		t.Fatal("WhenAnyAssetFailed should be false after load")
	}
}

func TestConditions_MultiAsset_AllLoadedAndAnyFailed(t *testing.T) {
	w, sch := initWorld(t)

	AddStaticFS(w, "assets", fstest.MapFS{
		"text/a.txt": &fstest.MapFile{Data: []byte("a")},
		// b.txt missing on purpose to force a failure
	})
	AddLoaderFn(w, testLoader, ".txt")

	AddAsset[testBundle2](w)
	b, ok := ecs.GetResource[testBundle2](w)
	if !ok {
		t.Fatal("testBundle2 resource not found after AddAsset")
	}
	if b.A == Asset(0) || b.B == Asset(0) {
		t.Fatal("expected both asset handles to be non-zero")
	}

	waitUntil(t, 2*time.Second,
		func() { tickAll(w, sch) },
		func() bool {
			// one should fail, one should load
			sa, oka := GetAssetState(w, b.A)
			sb, okb := GetAssetState(w, b.B)
			return oka && okb && (sa == AssetLoaded) && (sb == AssetFailed)
		},
	)

	if WhenAllAssetsLoaded(b.A, b.B)(w) {
		t.Fatal("WhenAllAssetsLoaded should be false when any asset failed")
	}
	if !WhenAnyAssetFailed(b.A, b.B)(w) {
		t.Fatal("WhenAnyAssetFailed should be true when at least one asset failed")
	}
}

func TestConditions_BundleWithZeroAssetStaysFalse(t *testing.T) {
	t.Parallel()

	w, _ := initWorld(t)

	ecs.AddResource(w, testBundle{Text: Asset(0)})

	if WhenBundleLoaded[testBundle]()(w) {
		t.Fatal("WhenBundleLoaded should be false when any Asset field is zero")
	}
}

func TestConcurrency_SetMaxConcurrentReads_StillLoadsManyAssets(t *testing.T) {
	w, sch := initWorld(t)

	m := fstest.MapFS{}
	const n = 25
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("text/%02d.txt", i)
		m[name] = &fstest.MapFile{Data: []byte(name)}
	}
	AddStaticFS(w, "assets", m)
	AddLoaderFn(w, testLoader, ".txt")

	type manyBundle struct {
		A0  Asset `path:"assets/text/00.txt"`
		A1  Asset `path:"assets/text/01.txt"`
		A2  Asset `path:"assets/text/02.txt"`
		A3  Asset `path:"assets/text/03.txt"`
		A4  Asset `path:"assets/text/04.txt"`
		A5  Asset `path:"assets/text/05.txt"`
		A6  Asset `path:"assets/text/06.txt"`
		A7  Asset `path:"assets/text/07.txt"`
		A8  Asset `path:"assets/text/08.txt"`
		A9  Asset `path:"assets/text/09.txt"`
		A10 Asset `path:"assets/text/10.txt"`
		A11 Asset `path:"assets/text/11.txt"`
		A12 Asset `path:"assets/text/12.txt"`
		A13 Asset `path:"assets/text/13.txt"`
		A14 Asset `path:"assets/text/14.txt"`
		A15 Asset `path:"assets/text/15.txt"`
		A16 Asset `path:"assets/text/16.txt"`
		A17 Asset `path:"assets/text/17.txt"`
		A18 Asset `path:"assets/text/18.txt"`
		A19 Asset `path:"assets/text/19.txt"`
		A20 Asset `path:"assets/text/20.txt"`
		A21 Asset `path:"assets/text/21.txt"`
		A22 Asset `path:"assets/text/22.txt"`
		A23 Asset `path:"assets/text/23.txt"`
		A24 Asset `path:"assets/text/24.txt"`
	}

	AddAsset[manyBundle](w)
	b, ok := ecs.GetResource[manyBundle](w)
	if !ok {
		t.Fatal("manyBundle resource not found after AddAsset")
	}

	assets := []Asset{
		b.A0, b.A1, b.A2, b.A3, b.A4, b.A5, b.A6, b.A7, b.A8, b.A9,
		b.A10, b.A11, b.A12, b.A13, b.A14, b.A15, b.A16, b.A17, b.A18, b.A19,
		b.A20, b.A21, b.A22, b.A23, b.A24,
	}

	waitUntil(t, 3*time.Second,
		func() { tickAll(w, sch) },
		func() bool { return WhenAllAssetsLoaded(assets...)(w) },
	)
}

func TestGetAndRemoveZeroAsset(t *testing.T) {
	t.Parallel()

	w, _ := initWorld(t)

	if v, ok := GetAsset[string](w, Asset(0)); ok || v != nil {
		t.Fatal("GetAsset with Asset(0) should return (nil,false)")
	}

	if ok := RemoveAsset[string](w, Asset(0)); ok {
		t.Fatal("RemoveAsset with Asset(0) should return false")
	}

	if _, ok := GetAssetState(w, Asset(0)); ok {
		t.Fatal("GetAssetState with Asset(0) should return ok=false")
	}
	if _, ok := GetAssetPath(w, Asset(0)); ok {
		t.Fatal("GetAssetPath with Asset(0) should return ok=false")
	}
	if _, ok := GetAssetError(w, Asset(0)); ok {
		t.Fatal("GetAssetError with Asset(0) should return ok=false")
	}
}
