package audio

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/vistormu/go-dsa/errors"
	"github.com/vistormu/xpeto/internal/core"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

const sampleRate = 48_000

type Library struct {
	nextId uint32

	context *Context

	refs   map[Audio]int
	paths  *core.BiHashmap[string, Audio]
	audios map[Audio]*Player
}

func NewLibrary() *Library {
	return &Library{
		context: NewContext(sampleRate),
		nextId:  1,
		refs:    make(map[Audio]int),
		paths:   core.NewBiHashmap[string, Audio](),
		audios:  make(map[Audio]*Player),
	}
}

func (l *Library) Register(path string) Audio {
	aud, ok := l.paths.GetByKey(path)
	if ok {
		return aud
	}

	aud = Audio{Id: l.nextId}
	l.nextId++
	l.paths.Put(path, aud)

	return aud
}

func (l *Library) Load(ctx *core.Context, audios ...Audio) {
	fsys, ok := core.GetResource[fs.FS](ctx)
	if !ok {
		errors.New(FilesystemNotSet).With("stage", "load").Print()
		return
	}

	for _, aud := range audios {
		// check if the audio is already loaded
		_, ok := l.audios[aud]
		if ok {
			l.refs[aud]++
			continue
		}

		// get the path from the aud
		path, ok := l.paths.GetByValue(aud)
		if !ok {
			errors.New(AudioHandleNotFound).With("audio", aud, "stage", "load").Print()
			continue
		}

		// load the audio
		audio, err := NewPlayer(fsys, path, l.context)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// store the audio in the map
		l.audios[aud] = audio
		l.refs[aud] = 1
	}
}

func (l *Library) Unload(audios ...Audio) {
	for _, aud := range audios {
		// check if the audio is loaded
		audio, ok := l.audios[aud]
		if !ok {
			errors.New(AudioHandleNotFound).With("audio", aud, "stage", "unload").Print()
			continue
		}

		// decrement the reference count
		l.refs[aud]--
		if l.refs[aud] > 0 {
			continue
		}

		// close the audio and remove it from the map
		audio.Close()
		delete(l.audios, aud)
		delete(l.refs, aud)
	}
}

func (m *Library) Player(aud Audio) (*Player, bool) {
	audio, ok := m.audios[aud]
	if !ok {
		// TODO: create default audio
		errors.New(AudioHandleNotFound).With("audio", aud, "stage", "get").Print()
		return nil, false
	}

	return audio, true
}

// =======
// helpers
// =======
func NewContext(sampleRate int) *Context {
	return audio.NewContext(sampleRate)
}

var decoders = map[string]func(io.Reader) (io.Reader, error){
	".wav": func(r io.Reader) (io.Reader, error) {
		return wav.DecodeWithSampleRate(sampleRate, r)
	},
	".mp3": func(r io.Reader) (io.Reader, error) {
		return mp3.DecodeWithSampleRate(sampleRate, r)
	},
}

func NewPlayer(fsys fs.FS, path string, context *Context) (*Player, error) {
	// read raw data from the file
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, errors.New(AudioPathNotFound).With("path", path).Wrap(err)
	}

	// check extension for decoder
	ext := strings.ToLower(filepath.Ext(path))

	_, ok := decoders[ext]
	if !ok {
		return nil, errors.New(UnsupportedAudioFormat).With("path", path, "extension", ext)
	}

	// decode the audio data
	streamer, err := decoders[ext](bytes.NewReader(data))
	if err != nil {
		return nil, errors.New(LoadAudioError).With("path", path, "extension", ext).Wrap(err)
	}

	// create a Player from the decoded stream
	audio, err := context.NewPlayer(streamer)
	if err != nil {
		return nil, errors.New(LoadAudioError).With("path", path, "extension", ext).Wrap(err)
	}

	return audio, nil
}
