package image

import (
	"fmt"
	"io/fs"

	"github.com/vistormu/xpeto/internal/errors"
	st "github.com/vistormu/xpeto/internal/structures"
)

type Manager struct {
	fsys   fs.FS
	nextId uint32

	refs   map[Handle]int
	paths  *st.BiHashmap[string, Handle]
	images map[Handle]*Image
}

func NewManager(fsys fs.FS) *Manager {
	return &Manager{
		fsys:   fsys,
		nextId: 0,
		refs:   make(map[Handle]int),
		paths:  st.NewBiHashmap[string, Handle](),
		images: make(map[Handle]*Image),
	}
}

func (m *Manager) Register(path string) Handle {
	// check if the path is already registered
	handle, ok := m.paths.GetByKey(path)
	if ok {
		return handle
	}

	handle = Handle{Id: m.nextId}
	m.nextId++
	m.paths.Put(path, handle)

	return handle
}

func (m *Manager) Load(handles ...Handle) {
	for _, handle := range handles {
		// check if the image is already loaded
		_, ok := m.images[handle]
		if ok {
			m.refs[handle]++
			continue
		}

		// get the path from the handle
		path, ok := m.paths.GetByValue(handle)
		if !ok {
			errors.New(errors.ImageHandleNotFound).With("handle", handle, "stage", "load").Print()
			continue
		}

		// load the image
		image, err := NewImage(m.fsys, path)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// store the image
		m.images[handle] = image
		m.refs[handle] = 1
	}
}

func (m *Manager) Unload(handles ...Handle) {
	for _, handle := range handles {
		// check if the image is loaded
		_, ok := m.images[handle]
		if !ok {
			errors.New(errors.ImageHandleNotFound).With("handle", handle, "stage", "unload").Print()
			continue
		}

		// decrement the reference count
		m.refs[handle]--
		if m.refs[handle] > 0 {
			continue
		}

		// remove the image from the map
		delete(m.images, handle)
		delete(m.refs, handle)
	}
}

func (m *Manager) Image(handle Handle) *Image {
	image, ok := m.images[handle]
	if !ok {
		// TODO: create default image
		errors.New(errors.ImageHandleNotFound).With("handle", handle, "stage", "get").Print()
		return nil
	}

	return image
}
