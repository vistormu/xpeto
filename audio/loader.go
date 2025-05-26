package audio

import (
	"bytes"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/vistormu/xpeto/internal/errors"
)

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

func NewAudio(fsys fs.FS, path string, context *Context) (*Audio, error) {
	// read raw data from the file
	data, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, errors.New(errors.AudioPathNotFound).With("path", path).Wrap(err)
	}

	// check extension for decoder
	ext := strings.ToLower(filepath.Ext(path))

	_, ok := decoders[ext]
	if !ok {
		return nil, errors.New(errors.UnsupportedAudioFormat).With("path", path, "extension", ext)
	}

	// decode the audio data
	streamer, err := decoders[ext](bytes.NewReader(data))
	if err != nil {
		return nil, errors.New(errors.LoadAudioError).With("path", path, "extension", ext).Wrap(err)
	}

	// create a Player from the decoded stream
	audio, err := context.NewPlayer(streamer)
	if err != nil {
		return nil, errors.New(errors.LoadAudioError).With("path", path, "extension", ext).Wrap(err)
	}

	return audio, nil
}
