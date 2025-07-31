package asset

import (
	"io"
	"io/fs"
	"reflect"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/event"
)

func TestServerRegistration(t *testing.T) {
	ctx := prepareMockContext()
	as := core.MustResource[*Server](ctx)

	// load the mock assets
	Load[txtType, mockAsset](as)

	// nextId
	if as.nextId != 3 {
		t.Errorf("Expected nextId to be 3, got %d", as.nextId)
	}

	// check if the assets are registered
	if as.registered.Length() != 2 {
		t.Errorf("Expected 2 registered assets, got %d", as.registered.Length())
	}

	// load states
	if len(as.loadStates) != 2 {
		t.Errorf("Expected 2 load states, got %d", len(as.loadStates))
	}

	// pending queue
	if as.pending.Length() != 2 {
		t.Errorf("Expected 2 pending assets, got %d", as.pending.Length())
	}
}

func TestAssetLoading(t *testing.T) {
	ctx := prepareMockContext()
	as := core.MustResource[*Server](ctx)
	eb := core.MustResource[*event.Bus](ctx)

	// event
	var wg sync.WaitGroup
	wg.Add(2)
	event.Subscribe(eb, func(data AssetEvent) {
		if data.Kind == Added {
			t.Log("data was successfully added!")
			wg.Done()
		}
	})

	// load the mock assets
	Load[txtType, mockAsset](as)

	// tick once
	Update(ctx)
	_, ok := core.GetResource[mockAsset](ctx)
	if !ok {
		t.Fatal("Expected mockAsset to be loaded, but it was not found in the context")
	}

	// mock game loop
	fps := time.NewTicker(time.Second / 60)
	defer fps.Stop()

	timeout := time.After(time.Second * 2)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

loop:
	for {
		select {
		case <-fps.C:
			Update(ctx)
		case <-done:
			t.Log("test successful")
			break loop
		case <-timeout:
			t.Fatalf("asset loading exceeded %v", timeout)
		}
	}

	// check that the assets exist in the context and that the state is Loaded
	mockAssets, _ := core.GetResource[mockAsset](ctx)

	t.Log(mockAssets)

	state1 := as.GetState(mockAssets.Asset1)
	state2 := as.GetState(mockAssets.Asset2)

	if state1 != Loaded || state2 != Loaded {
		t.Fatalf("assets do not appear as loaded!: %d, %d", state1, state2)
	}

	t.Log(as.assetStore[reflect.TypeFor[txtType]()])
	t.Log(mockAssets.Asset1)
	t.Log(mockAssets.Asset2)

	// check that the assets are loaded
	asset1, ok := GetAsset[txtType](as, mockAssets.Asset1)
	if !ok {
		t.Fatal("Expected asset1 to be loaded, but it was not found in the Asset")
	}

	asset2, ok := GetAsset[txtType](as, mockAssets.Asset2)
	if !ok {
		t.Fatal("Expected asset2 to be loaded, but it was not found in the Asset")
	}

	if string(asset1) != "This is asset 1" {
		t.Errorf("Expected asset1 content to be 'This is asset 1', got '%s'", asset1)
	}

	if string(asset2) != "This is asset 2" {
		t.Errorf("Expected asset2 content to be 'This is asset 2', got '%s'", asset2)
	}
}

// =======
// helpers
// =======
type mockAsset struct {
	Asset1 Handle `path:"assets/asset1.txt"`
	Asset2 Handle `path:"assets/asset2.txt"`
}

type txtType string

func createMockAssets() fs.FS {
	mockFS := fstest.MapFS{
		"assets/asset1.txt": &fstest.MapFile{
			Data: []byte("This is asset 1"),
			Mode: fs.ModePerm,
		},
		"assets/asset2.txt": &fstest.MapFile{
			Data: []byte("This is asset 2"),
			Mode: fs.ModePerm,
		},
	}
	return mockFS
}

func mockLoader(reader io.Reader, path string) (any, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return txtType(data), nil
}

func prepareMockContext() *core.Context {
	ctx := core.NewContext()
	server := NewServer()
	fsys := createMockAssets()

	server.AddLoader(".txt", mockLoader)
	server.SetFilesystem(fsys)

	core.AddResource(ctx, server)

	eb := event.NewBus()
	core.AddResource(ctx, eb)

	return ctx
}
