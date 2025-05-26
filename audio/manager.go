package audio

import (
	"fmt"
	"io/fs"

	"github.com/vistormu/xpeto/internal/errors"
	st "github.com/vistormu/xpeto/internal/structures"
)

const sampleRate = 48_000

type Manager struct {
	fsys   fs.FS
	nextId uint32

	context *Context

	refs   map[Handle]int
	paths  *st.BiHashmap[string, Handle]
	audios map[Handle]*Audio
}

func NewManager(fsys fs.FS) *Manager {
	return &Manager{
		fsys:    fsys,
		context: NewContext(sampleRate),
		nextId:  0,
		refs:    make(map[Handle]int),
		paths:   st.NewBiHashmap[string, Handle](),
		audios:  make(map[Handle]*Audio),
	}
}

func (m *Manager) Register(path string) Handle {
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
		// check if the audio is already loaded
		_, ok := m.audios[handle]
		if ok {
			m.refs[handle]++
			continue
		}

		// get the path from the handle
		path, ok := m.paths.GetByValue(handle)
		if !ok {
			errors.New(errors.AudioHandleNotFound).With("handle", handle, "stage", "load").Print()
			continue
		}

		// load the audio
		audio, err := NewAudio(m.fsys, path, m.context)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// store the audio in the map
		m.audios[handle] = audio
		m.refs[handle] = 1
	}
}

func (m *Manager) Unload(handles ...Handle) {
	for _, handle := range handles {
		// check if the audio is loaded
		audio, ok := m.audios[handle]
		if !ok {
			errors.New(errors.AudioHandleNotFound).With("handle", handle, "stage", "unload").Print()
			continue
		}

		// decrement the reference count
		m.refs[handle]--
		if m.refs[handle] > 0 {
			continue
		}

		// close the audio and remove it from the map
		audio.Close()
		delete(m.audios, handle)
		delete(m.refs, handle)
	}
}

func (m *Manager) Audio(handle Handle) (*Audio, bool) {
	audio, ok := m.audios[handle]
	if !ok {
		// TODO: create default audio
		errors.New(errors.AudioHandleNotFound).With("handle", handle, "stage", "get").Print()
		return nil, false
	}

	return audio, true
}
