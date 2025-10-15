package asset

import (
	"io/fs"
	"log"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/vistormu/go-dsa/hashmap"
	"github.com/vistormu/go-dsa/queue"

	"github.com/vistormu/xpeto/core/ecs"
)

type loadResult struct {
	req     loadRequest
	content any
	err     error
}

type Server struct {
	mu   sync.RWMutex
	fsys fs.FS

	nextId      int
	freeHandles *queue.QueueArray[Handle]

	registered *hashmap.BiHashmap[string, Handle]
	loadStates map[string]LoadState

	pending   *queue.QueueArray[loadRequest]
	completed chan loadResult

	loaders    map[string]loaderFn
	assetStore map[reflect.Type]map[Handle]any
}

func NewServer() *Server {
	return &Server{
		mu:          sync.RWMutex{},
		fsys:        nil,
		nextId:      1,
		freeHandles: queue.NewQueueArray[Handle](),
		registered:  hashmap.NewBiHashmap[string, Handle](),
		loadStates:  make(map[string]LoadState),
		pending:     queue.NewQueueArray[loadRequest](),
		completed:   make(chan loadResult, 16),
		loaders:     make(map[string]loaderFn),
		assetStore:  make(map[reflect.Type]map[Handle]any),
	}
}

func storeAssetByType(s *Server, type_ reflect.Type, handle Handle, asset any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// check if the asset store exists for the type
	assets, ok := s.assetStore[type_]
	if !ok {
		assets = make(map[Handle]any)
		s.assetStore[type_] = assets
	}

	// store the asset with the handle
	assets[handle] = asset

	// update the load state to loaded
	path, ok := s.registered.GetByValue(handle)
	if !ok {
		log.Printf("handle %v not found in registered assets", handle)
		return
	}
	s.loadStates[path] = Loaded
}

// ===
// API
// ===
func SetFileSystem(w *ecs.World, fsys fs.FS) {
	as, _ := ecs.GetResource[Server](w)
	as.fsys = fsys
}

func AddAssetLoader(w *ecs.World, extension string, loader loaderFn) {
	as, _ := ecs.GetResource[Server](w)

	as.mu.Lock()
	defer as.mu.Unlock()

	_, ok := as.loaders[extension]
	if ok {
		return
	}

	as.loaders[extension] = loader
}

func AddAssets[T any, B any](w *ecs.World) {
	as, _ := ecs.GetResource[Server](w)

	bundle := new(B)
	bundleValue := reflect.ValueOf(bundle).Elem()
	bundleType := reflect.TypeFor[B]()
	assetType := reflect.TypeFor[T]()

	for i := 0; i < bundleValue.NumField(); i++ {
		// get path from tag
		path := bundleType.Field(i).Tag.Get("path")
		if path == "" {
			panic("missing path tag for field " + bundleType.Field(i).Name)
		}

		// check if the path is already registered
		_, ok := as.registered.GetByKey(path)
		if ok {
			continue
		}

		// get a free id or create a new one
		var handle Handle
		handle, err := as.freeHandles.Dequeue()
		if err == nil {
			handle = Handle{
				Number:  handle.Number,
				Version: handle.Version + 1,
			}
		} else {
			handle = Handle{
				Number:  as.nextId,
				Version: 1,
			}
			as.nextId++
		}

		// set the bundleValue of the handle
		field := bundleValue.Field(i)
		if field.CanSet() && field.Type() == reflect.TypeOf(Handle{}) {
			field.Set(reflect.ValueOf(handle))
		} else {
			log.Printf("field %s must be of type Handle and settable", bundleType.Field(i).Name)
			continue
		}

		// register the path with the handle
		as.registered.Put(path, handle)

		// set new load state
		as.loadStates[path] = Loading

		// create the load context and put it in the pending queue
		loader, ok := as.loaders[filepath.Ext(path)]
		if !ok {
			log.Println("loader not found")
			continue
		}
		as.pending.Enqueue(loadRequest{
			path:       path,
			bundle:     *bundle,
			bundleType: bundleType,
			handle:     handle,
			assetType:  assetType,
			loaderFn:   loader,
		})
	}
}

func GetAsset[T any](w *ecs.World, handle Handle) (T, bool) {
	as, _ := ecs.GetResource[Server](w)

	as.mu.RLock()
	defer as.mu.RUnlock()

	assetType := reflect.TypeFor[T]()

	assets, ok := as.assetStore[assetType]
	if !ok {
		var zero T
		return zero, false
	}

	asset, ok := assets[handle]
	if !ok {
		var zero T
		return zero, false
	}

	out, ok := asset.(T)
	if !ok {
		var zero T
		return zero, false
	}

	return out, true
}

// func StoreAsset[T any](s *Server, handle Handle, asset T) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	assetType := reflect.TypeFor[T]()

// 	// check if the asset store exists for the type
// 	assets, ok := s.assetStore[assetType]
// 	if !ok {
// 		assets = make(map[Handle]any)
// 		s.assetStore[assetType] = assets
// 	}

// 	// store the asset with the handle
// 	assets[handle] = asset

// 	// update the load state to loaded
// 	path, ok := s.registered.GetByValue(handle)
// 	if !ok {
// 		log.Printf("handle %v not found in registered assets", handle)
// 		return
// 	}
// 	s.loadStates[path] = Loaded
// }

// func (s *Server) GetState(handle Handle) LoadState {
// 	path, ok := s.registered.GetByValue(handle)
// 	if !ok {
// 		return NotFound
// 	}

// 	state, ok := s.loadStates[path]
// 	if !ok {
// 		return NotFound
// 	}

// 	return state
// }
