package text

import (
	"io/fs"
)

type Manager struct{}

func NewManager(fsys fs.FS) *Manager {
	return &Manager{}
}
