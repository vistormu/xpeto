package asset

import (
	"testing"
	"testing/fstest"
	"time"

	"github.com/vistormu/xpeto/core"
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/schedule"
)

func initWorld() *ecs.World {
	w := ecs.NewWorld()
	sch := schedule.NewScheduler()

	core.CorePkgs(w, sch)

	return w
}

func TestSplitPath(t *testing.T) {
	tests := []struct {
		in       string
		wantBase string
		wantRel  string
		wantExt  string
	}{
		{
			in:       "assets/images/logo.png",
			wantBase: "assets",
			wantRel:  "images/logo.png",
			wantExt:  ".png",
		},
		{
			in:       "audio/sfx/click.WAV",
			wantBase: "audio",
			wantRel:  "sfx/click.WAV",
			wantExt:  ".wav",
		},
		{
			in:       "file.txt",
			wantBase: "",
			wantRel:  "file.txt",
			wantExt:  ".txt",
		},
		{
			in:       "noext",
			wantBase: "",
			wantRel:  "noext",
			wantExt:  "",
		},
	}

	for _, tt := range tests {
		base, rel, ext := splitPath(tt.in)
		if base != tt.wantBase || rel != tt.wantRel || ext != tt.wantExt {
			t.Fatalf("splitPath(%q) = (%q,%q,%q), want (%q,%q,%q)",
				tt.in, base, rel, ext, tt.wantBase, tt.wantRel, tt.wantExt)
		}
	}
}

func TestAddStaticFSSanitiseBase(t *testing.T) {
	w := initWorld()
	ecs.AddResource(w, newServer())

	// empty base -> no registration
	AddStaticFS(w, "", fstest.MapFS{})
	s, _ := ecs.GetResource[server](w)
	if len(s.staticFS) != 0 {
		t.Fatalf("expected no filesystems registered for empty base, got %d", len(s.staticFS))
	}

	// base with path separators -> no registration
	AddStaticFS(w, "bad/name", fstest.MapFS{})
	if len(s.staticFS) != 0 {
		t.Fatalf("expected no filesystems registered for base with slash, got %d", len(s.staticFS))
	}

	// valid base
	AddStaticFS(w, "assets", fstest.MapFS{})
	if len(s.staticFS) != 1 {
		t.Fatalf("expected 1 filesystem registered, got %d", len(s.staticFS))
	}
	if _, ok := s.staticFS["assets"]; !ok {
		t.Fatalf("filesystem with base %q not found", "assets")
	}
}

type testBundle struct {
	Text Asset `path:"assets/text/hello.txt"`
}

type testType struct {
	Value string
}

func testLoader(data []byte, path string) (*testType, error) {
	return &testType{
		Value: string(data),
	}, nil
}

func TestAssetLoadEndToEnd(t *testing.T) {
	w := initWorld()
	ecs.AddResource(w, newServer())
	ecs.AddResource(w, newLoader())

	// filesystem with one text asset
	fsys := fstest.MapFS{
		"text/hello.txt": &fstest.MapFile{
			Data: []byte("hello, world"),
			Mode: 0o444,
		},
	}

	// register filesystem and loader
	AddStaticFS(w, "assets", fsys)
	AddLoaderFn(w, testLoader, ".txt")

	// register bundle
	AddAsset[testBundle](w)

	// retrieve bundle resource
	bundle, ok := ecs.GetResource[testBundle](w)
	if !ok {
		t.Fatal("testBundle resource not found after AddAsset")
	}
	if bundle.Text == Asset(0) {
		t.Fatal("asset handle for Text is zero, expected non-zero")
	}

	const (
		maxTries = 50
		delay    = 10 * time.Millisecond
	)

	var got *testType
	var times int
	for range maxTries {
		readRequests(w)
		loadResults(w)

		v, ok := GetAsset[testType](w, bundle.Text)
		if ok {
			got = v
			break
		}
		time.Sleep(delay)

		times++
	}

	t.Logf("loaded after %d ticks", times)

	if got == nil {
		t.Fatal("asset not loaded in time")
	}
	if got.Value != "hello, world" {
		t.Fatalf("loaded asset value = %q, want %q", *got, "hello, world")
	}

	s, _ := ecs.GetResource[server](w)

	if s.store.Len() != 1 {
		t.Fatalf("there should be only 1 store and got %d", s.store.Len())
	}

	store := getStore[testType](s.store)
	if len(store.dense) != 1 {
		t.Fatalf("there should be only one asset and got %d", len(store.dense))
	}

	if len(store.values) != 1 {
		t.Fatalf("there should be only one asset and got %d", len(store.values))
	}

	// test RemoveAsset
	if !RemoveAsset[testType](w, bundle.Text) {
		t.Fatal("RemoveAsset returned false, expected true")
	}
	if v, ok := GetAsset[testType](w, bundle.Text); ok || v != nil {
		t.Fatal("expected asset to be removed, but GetAsset still returned a value")
	}
}

func TestGetAndRemoveZeroAsset(t *testing.T) {
	w := initWorld()
	ecs.AddResource(w, newServer())

	// GetAsset with zero handle
	if v, ok := GetAsset[string](w, Asset(0)); ok || v != nil {
		t.Fatal("GetAsset with Asset(0) should return (nil,false)")
	}

	// RemoveAsset with zero handle
	if ok := RemoveAsset[string](w, Asset(0)); ok {
		t.Fatal("RemoveAsset with Asset(0) should return false")
	}
}
