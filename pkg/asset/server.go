package asset

import (
	"io/fs"
	"log"
	"path/filepath"
	"reflect"
	"sync"

	"github.com/vistormu/xpeto/internal/core"
)

type loadResult struct {
	req     LoadRequest
	content any
	err     error
}

type Server struct {
	mu   sync.RWMutex
	fsys fs.FS

	nextId      uint32
	freeHandles *core.QueueArray[Handle]

	registered *core.BiHashmap[string, Handle]
	loadStates map[string]LoadState

	pending   *core.QueueArray[LoadRequest]
	completed chan loadResult

	loaders    map[string]LoaderFn
	assetStore map[reflect.Type]map[Handle]any
}

func NewServer() *Server {
	return &Server{
		mu:          sync.RWMutex{},
		fsys:        nil,
		nextId:      1,
		freeHandles: core.NewQueueArray[Handle](),
		registered:  core.NewBiHashmap[string, Handle](),
		loadStates:  make(map[string]LoadState),
		pending:     core.NewQueueArray[LoadRequest](),
		completed:   make(chan loadResult, 16),
		loaders:     make(map[string]LoaderFn),
		assetStore:  make(map[reflect.Type]map[Handle]any),
	}
}

func (s *Server) SetFilesystem(fsys fs.FS) {
	s.fsys = fsys
}

// =======
// loaders
// =======
func (s *Server) AddLoader(extension string, loader LoaderFn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.loaders[extension]
	if ok {
		return
	}

	s.loaders[extension] = loader
}

func (s *Server) GetLoader(extension string) (LoaderFn, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	loader, ok := s.loaders[extension]
	return loader, ok
}

// ======
// assets
// ======
func Load[T any, B any](s *Server) {
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
		_, ok := s.registered.GetByKey(path)
		if ok {
			continue
		}

		// get a free id or create a new one
		var handle Handle
		handle, err := s.freeHandles.Dequeue()
		if err == nil {
			handle = Handle{
				Id:      handle.Id,
				Version: handle.Version + 1,
			}
		} else {
			handle = Handle{
				Id:      s.nextId,
				Version: 1,
			}
			s.nextId++
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
		s.registered.Put(path, handle)

		// set new load state
		s.loadStates[path] = Loading

		// create the load context and put it in the pending queue
		loader, ok := s.loaders[filepath.Ext(path)]
		if !ok {
			log.Println("loader not found")
			continue
		}
		s.pending.Enqueue(LoadRequest{
			Path:       path,
			Bundle:     *bundle,
			BundleType: bundleType,
			Handle:     handle,
			AssetType:  assetType,
			LoaderFn:   loader,
		})
	}
}

func StoreAsset[T any](s *Server, handle Handle, asset T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	assetType := reflect.TypeFor[T]()

	// check if the asset store exists for the type
	assets, ok := s.assetStore[assetType]
	if !ok {
		assets = make(map[Handle]any)
		s.assetStore[assetType] = assets
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

func StoreAssetByType(s *Server, type_ reflect.Type, handle Handle, asset any) {
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

func GetAsset[T any](s *Server, handle Handle) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	assetType := reflect.TypeFor[T]()

	assets, ok := s.assetStore[assetType]
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

// =====
// state
// =====
func (s *Server) GetState(handle Handle) LoadState {
	path, ok := s.registered.GetByValue(handle)
	if !ok {
		return NotFound
	}

	state, ok := s.loadStates[path]
	if !ok {
		return NotFound
	}

	return state
}
