package font

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type font struct {
	*text.GoTextFaceSource

	mu    sync.RWMutex
	faces map[float64]text.Face
}

func (f *font) Face(size float64) text.Face {
	f.mu.RLock()
	if f.faces != nil {
		if face, ok := f.faces[size]; ok {
			f.mu.RUnlock()
			return face
		}
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()

	if f.faces == nil {
		f.faces = make(map[float64]text.Face)
	}
	if face, ok := f.faces[size]; ok {
		return face
	}
	face := &text.GoTextFace{Source: f.GoTextFaceSource, Size: size}
	f.faces[size] = face

	return face
}
