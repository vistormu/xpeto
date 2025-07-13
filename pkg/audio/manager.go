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

	refs   map[Audio]int
	paths  *st.BiHashmap[string, Audio]
	audios map[Audio]*Player
}

func NewManager() *Manager {
	return &Manager{
		fsys:    nil,
		context: NewContext(sampleRate),
		nextId:  1,
		refs:    make(map[Audio]int),
		paths:   st.NewBiHashmap[string, Audio](),
		audios:  make(map[Audio]*Player),
	}
}

func (m *Manager) WithFilesystem(fsys fs.FS) *Manager {
	m.fsys = fsys
	return m
}

func (m *Manager) Register(path string) Audio {
	aud, ok := m.paths.GetByKey(path)
	if ok {
		return aud
	}

	aud = Audio{Id: m.nextId}
	m.nextId++
	m.paths.Put(path, aud)

	return aud
}

func (m *Manager) Load(audios ...Audio) {
	if m.fsys == nil {
		errors.New(errors.FilesystemNotSet).With("stage", "load").Print()
		return
	}

	for _, aud := range audios {
		// check if the audio is already loaded
		_, ok := m.audios[aud]
		if ok {
			m.refs[aud]++
			continue
		}

		// get the path from the aud
		path, ok := m.paths.GetByValue(aud)
		if !ok {
			errors.New(errors.AudioHandleNotFound).With("audio", aud, "stage", "load").Print()
			continue
		}

		// load the audio
		audio, err := NewPlayer(m.fsys, path, m.context)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// store the audio in the map
		m.audios[aud] = audio
		m.refs[aud] = 1
	}
}

func (m *Manager) Unload(audios ...Audio) {
	if m.fsys == nil {
		errors.New(errors.FilesystemNotSet).With("stage", "unload").Print()
		return
	}

	for _, aud := range audios {
		// check if the audio is loaded
		audio, ok := m.audios[aud]
		if !ok {
			errors.New(errors.AudioHandleNotFound).With("audio", aud, "stage", "unload").Print()
			continue
		}

		// decrement the reference count
		m.refs[aud]--
		if m.refs[aud] > 0 {
			continue
		}

		// close the audio and remove it from the map
		audio.Close()
		delete(m.audios, aud)
		delete(m.refs, aud)
	}
}

func (m *Manager) Player(aud Audio) (*Player, bool) {
	audio, ok := m.audios[aud]
	if !ok {
		// TODO: create default audio
		errors.New(errors.AudioHandleNotFound).With("audio", aud, "stage", "get").Print()
		return nil, false
	}

	return audio, true
}
