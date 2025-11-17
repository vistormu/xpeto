package window

import (
	"github.com/vistormu/go-dsa/constraints"

	"github.com/vistormu/xpeto/core/ecs"
)

// ====
// real
// ====
type RealWindow struct {
	Title        string
	Width        int
	Height       int
	FullScreen   bool
	AntiAliasing bool
}

type lastRealWindow struct {
	RealWindow
}

// =======
// virtual
// =======
type VirtualWindow struct {
	Width  int
	Height int
}

// ===
// API
// ===
func SetRealWindowSize(w *ecs.World, width, height int) {
	rw, _ := ecs.GetResource[RealWindow](w)
	rw.Width = width
	rw.Height = height
}

func GetRealWindowSize[T constraints.Number](w *ecs.World) (width, height T) {
	rw, _ := ecs.GetResource[RealWindow](w)
	width = T(rw.Width)
	height = T(rw.Height)

	return
}

func SetVirtualWindowSize(w *ecs.World, width, height int) {
	rw, _ := ecs.GetResource[VirtualWindow](w)
	rw.Width = width
	rw.Height = height
}

func GetVirtualWindowSize[T constraints.Number](w *ecs.World) (width, height T) {
	rw, _ := ecs.GetResource[VirtualWindow](w)
	width = T(rw.Width)
	height = T(rw.Height)

	return
}
