package render

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vistormu/xpeto/internal/core"
)

type Renderable interface {
	SortKey() uint64
	Draw(screen *ebiten.Image)
}

type ExtractionFn func(*core.Context) []Renderable

type Phase int

const (
	Transparent Phase = iota
	Opaque
	UI
	PostFx
)
