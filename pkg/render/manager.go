package render

import (
	"fmt"
	"io/fs"

	"github.com/vistormu/xpeto/internal/errors"
	st "github.com/vistormu/xpeto/internal/structures"
)

type Manager struct {
	fsys   fs.FS
	nextId uint32

	refs      map[Image]int
	paths     *st.BiHashmap[string, Image]
	renderers map[Image]*Renderer
}

func NewManager() *Manager {
	return &Manager{
		fsys:      nil,
		nextId:    1,
		refs:      make(map[Image]int),
		paths:     st.NewBiHashmap[string, Image](),
		renderers: make(map[Image]*Renderer),
	}
}

func (m *Manager) WithFilesystem(fsys fs.FS) *Manager {
	m.fsys = fsys
	return m
}

func (m *Manager) Register(path string) Image {
	// check if the path is already registered
	img, ok := m.paths.GetByKey(path)
	if ok {
		return img
	}

	img = Image{Id: m.nextId}
	m.nextId++
	m.paths.Put(path, img)

	return img
}

func (m *Manager) Load(images ...Image) {
	if m.fsys == nil {
		errors.New(errors.FilesystemNotSet).With("stage", "load").Print()
		return
	}

	for _, img := range images {
		// check if the image is already loaded
		_, ok := m.renderers[img]
		if ok {
			m.refs[img]++
			continue
		}

		// get the path from the img
		path, ok := m.paths.GetByValue(img)
		if !ok {
			errors.New(errors.ImageHandleNotFound).With("img", img, "stage", "load").Print()
			continue
		}

		// load the image
		renderer, err := NewRenderer(m.fsys, path)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// store the image
		m.renderers[img] = renderer
		m.refs[img] = 1
	}
}

func (m *Manager) Unload(images ...Image) {
	if m.fsys == nil {
		errors.New(errors.FilesystemNotSet).With("stage", "unload").Print()
		return
	}

	for _, img := range images {
		// check if the image is loaded
		_, ok := m.renderers[img]
		if !ok {
			errors.New(errors.ImageHandleNotFound).With("img", img, "stage", "unload").Print()
			continue
		}

		// decrement the reference count
		m.refs[img]--
		if m.refs[img] > 0 {
			continue
		}

		// remove the image from the map
		delete(m.renderers, img)
		delete(m.refs, img)
	}
}

func (m *Manager) Renderer(img Image) *Renderer {
	renderer, ok := m.renderers[img]
	if !ok {
		// TODO: create default image
		errors.New(errors.ImageHandleNotFound).With("img", img, "stage", "get").Print()
		return nil
	}

	return renderer
}
