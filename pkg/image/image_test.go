package image

import (
	"io/fs"
	"os"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"
	"github.com/vistormu/xpeto/internal/event"
	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/transform"
)

func TestLoadImage(t *testing.T) {
	ctx := prepareMockContext()
	as := core.MustResource[*asset.Server](ctx)
	eb := core.MustResource[*event.Bus](ctx)
	w := core.MustResource[*ecs.World](ctx)

	// load the mock assets
	asset.Load[*Image, mockImageBundle](as)

	// event
	var wg sync.WaitGroup
	wg.Add(1)
	event.Subscribe(eb, func(data asset.AssetEvent) {
		if data.Kind == asset.Added {
			t.Log("data was successfully added!")
			wg.Done()
		}
	})

	asset.Update(ctx)

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
			asset.Update(ctx)
		case <-done:
			t.Log("test successful")
			break loop
		case <-timeout:
			t.Fatalf("asset loading exceeded %v", timeout)
		}
	}

	mockBundle, ok := core.GetResource[mockImageBundle](ctx)
	if !ok {
		t.Fatal("Expected mockImageBundle to be loaded, but it was not found in the context")
	}

	state := as.GetState(mockBundle.DefaultImage)

	if state != asset.Loaded {
		t.Fatalf("Expected asset state to be Loaded, but got %v", state)
	}

	image, ok := asset.GetAsset[*Image](as, mockBundle.DefaultImage)
	if !ok {
		t.Fatal("Expected Image asset to be loaded, but it was not found in the context")
	}

	if image == nil {
		t.Fatal("Expected Image asset to be non-nil, but it was nil")
	}

	// add the renderable component to the world
	entity := w.Create()

	ecs.AddComponent(w, entity, &Renderable{
		Image: mockBundle.DefaultImage,
		Layer: 0,
	})

	ecs.AddComponent(w, entity, &transform.Transform{
		Position: core.Vector[float32]{X: 100, Y: 100},
		Scale:    core.Vector[float32]{X: 1, Y: 1},
		Rotation: 0,
	})
}

// =======
// helpers
// =======
type mockImageBundle struct {
	DefaultImage asset.Handle `path:"assets/default.png"`
}

func createMockAssets() fs.FS {
	data, err := os.ReadFile("../../assets/default.png")
	if err != nil {
		panic("failed to read mock asset file: " + err.Error())
	}

	mockFS := fstest.MapFS{
		"assets/default.png": &fstest.MapFile{
			Data: data,
			Mode: fs.ModePerm,
		},
	}

	return mockFS
}

func prepareMockContext() *core.Context {
	ctx := core.NewContext()
	fsys := createMockAssets()

	// asset server
	server := asset.NewServer()
	server.AddLoader(".png", LoadImage)
	server.SetFilesystem(fsys)
	core.AddResource(ctx, server)

	// event bus
	eb := event.NewBus()
	core.AddResource(ctx, eb)

	// world
	world := ecs.NewWorld()
	core.AddResource(ctx, world)

	return ctx
}
