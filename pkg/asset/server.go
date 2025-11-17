package asset

import (
	"io/fs"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/vistormu/go-dsa/hashmap"

	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/pkg/log"
)

// ======
// server
// ======
type server struct {
	population *population
	store      *hashmap.TypeMap

	staticFS map[string]fs.FS
	loaders  map[string]loaderHandler
}

func newServer() server {
	return server{
		population: newPopulation(),
		store:      hashmap.NewTypeMap(),
		staticFS:   make(map[string]fs.FS),
		loaders:    make(map[string]loaderHandler),
	}
}

// =======
// helpers
// =======
func baseType(t reflect.Type) reflect.Type {
	if t == nil {
		return nil
	}
	if t.Kind() == reflect.Pointer {
		return t.Elem()
	}
	return t
}

func splitPath(path string) (base string, rel string, ext string) {
	clean := filepath.ToSlash(filepath.Clean(path))

	ext = strings.ToLower(filepath.Ext(clean))

	parts := strings.SplitN(clean, "/", 2)
	if len(parts) != 2 {
		base = ""
		rel = clean
		return
	}

	base = parts[0]
	rel = parts[1]

	return
}

// ===
// API
// ===
func AddStaticFS(w *ecs.World, base string, fsys fs.FS) {
	s, ok := ecs.GetResource[server](w)
	if !ok {
		log.LogError(w, "cannot execute AddStaticFS: asset.Pkg not included")
		return
	}

	base = strings.TrimSpace(base)
	if base == "" {
		log.LogError(w, "cannot register filesystem: base is empty")
		return
	}

	if strings.Contains(base, "\\") || strings.Contains(base, "/") {
		log.LogError(w, "cannot register filesystem with base name containing path separators", log.F("base", base))
		return
	}

	_, ok = s.staticFS[base]
	if ok {
		log.LogWarning(w, "base name for filesystem already exists, overwriting it", log.F("name", base))
	}

	s.staticFS[base] = fsys
}

func AddLoaderFn[T any](w *ecs.World, fn LoaderFn[T], extensions ...string) {
	s, ok := ecs.GetResource[server](w)
	if !ok {
		log.LogError(w, "cannot execute AddLoaderFn: asset.Pkg not included")
		return
	}

	if len(extensions) == 0 {
		log.LogError(w, "AddLoaderFn called without extensions")
		return
	}

	for _, ext := range extensions {
		ext = strings.ToLower(ext)
		if ext == "" {
			log.LogWarning(w, "empty extension provided in AddLoaderFn")
			continue
		}
		if ext[0] != '.' {
			ext = "." + ext
		}

		_, ok := s.loaders[ext]
		if ok {
			log.LogWarning(w, "overwritting loader for extension", log.F("ext", ext))
		}

		s.loaders[ext] = func(s *server, req request, data []byte) error {
			v, err := fn(data, req.path)
			if err != nil {
				return err
			}

			getStore[T](s.store).add(req.asset, v)

			return nil
		}
	}
}

func AddAsset[T any](w *ecs.World) {
	s, ok := ecs.GetResource[server](w)
	if !ok {
		log.LogError(w, "cannot execute AddAsset: asset.Pkg not included")
		return
	}
	l, _ := ecs.GetResource[loader](w)

	// types
	b := new(T)
	bType := baseType(reflect.TypeFor[T]())
	bValue := reflect.ValueOf(b).Elem()

	if bType.Kind() != reflect.Struct {
		log.LogError(w, "the asset bundle must be a struct", log.F("got", bType.Kind().String()))
	}

	// iterate over the bundle
	for i := range bValue.NumField() {
		// get path from tag
		path := bType.Field(i).Tag.Get("path")
		if path == "" {
			log.LogError(w, "missing path tag for field", log.F("name", bType.Field(i).Name))
			continue
		}

		base, _, ext := splitPath(path)
		_, ok := s.staticFS[base]
		if !ok {
			log.LogError(w, "missing filesystem for asset", log.F("base", base))
			continue
		}

		if ext == "" {
			log.LogError(w, "missing extension on path for asset", log.F("path", path))
			continue
		}

		_, ok = s.loaders[ext]
		if !ok {
			log.LogError(w, "missing loader for asset type", log.F("ext", ext))
			continue
		}

		// check field type
		field := bValue.Field(i)
		if !field.CanSet() {
			log.LogError(w, "field must be settable", log.F("field", bType.Field(i).Name))
			continue
		}

		if field.Type() != reflect.TypeFor[Asset]() {
			log.LogError(w, "field must be of type Asset", log.F("got", field.Type().String()))
			continue
		}

		// add asset
		a := s.population.add()
		field.Set(reflect.ValueOf(a))

		// add request
		l.enqueue(request{
			path:  path,
			asset: a,
		})
	}

	// add bundle to the resources
	ecs.AddResource(w, b)
}

func GetAsset[T any](w *ecs.World, a Asset) (*T, bool) {
	s, ok := ecs.GetResource[server](w)
	if !ok {
		log.LogError(w, "cannot execute GetAsset: asset.Pkg not included")
		return nil, false
	}

	if a == Asset(0) {
		log.LogWarning(w, "cannot get asset with an initialized Asset type")
		return nil, false
	}

	return getStore[T](s.store).get(a)
}

func RemoveAsset[T any](w *ecs.World, a Asset) bool {
	s, ok := ecs.GetResource[server](w)
	if !ok {
		log.LogError(w, "cannot execute RemoveAsset: asset.Pkg not included")
		return false
	}

	if a == Asset(0) {
		log.LogWarning(w, "cannot remove asset with an initialized Asset type")
		return false
	}

	s.population.remove(a)

	return getStore[T](s.store).remove(a)
}
