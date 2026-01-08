package font

import (
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type font struct {
	*text.GoTextFaceSource

	mu    sync.RWMutex
	faces map[float64]text.Face
}

const (
	minFaceSize  = 4.0
	maxFaceSize  = 256.0
	faceSizeStep = 0.25
)

func sanitizeFaceSize(size float64) float64 {
	if math.IsNaN(size) || math.IsInf(size, 0) {
		return 18
	}

	size = max(minFaceSize, min(size, maxFaceSize))
	steps := math.Round(size / faceSizeStep)

	return steps * faceSizeStep
}

func (f *font) Face(size float64) text.Face {
	size = sanitizeFaceSize(size)

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
