package window

import (
	"github.com/vistormu/go-dsa/constraints"

	"github.com/vistormu/xpeto/core/ecs"
)

// ====
// real
// ====
type ResizingMode uint8

const (
	ResizingModeDisabled ResizingMode = iota
	ResizingModeOnlyFullscreenEnabled
	ResizingModeEnabled
)

type SizeLimits struct {
	MinW, MinH int
	MaxW, MaxH int
}

type WindowAction uint8

const (
	ActionNone WindowAction = iota
	ActionMaximize
	ActionMinimize
	ActionRestore
)

type RealWindow struct {
	Title               string
	Width               int
	Height              int
	FullScreen          bool
	AntiAliasing        bool
	VSync               bool
	RunnableOnUnfocused bool
	ResizingMode        ResizingMode
	SizeLimits          SizeLimits
	Action              WindowAction
}

type RealWindowObserved struct {
	Width       int
	Height      int
	DeviceScale float64
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

// real
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

func SetFullScreen(w *ecs.World, v bool) {
	rw, _ := ecs.GetResource[RealWindow](w)
	rw.FullScreen = v
}

func SetAntiAliasing(w *ecs.World, v bool) {
	rw, _ := ecs.GetResource[RealWindow](w)
	rw.AntiAliasing = v
}

func SetVSync(w *ecs.World, v bool) {
	rw, _ := ecs.GetResource[RealWindow](w)
	rw.VSync = v
}

func SetRunnableOnUnfocused(w *ecs.World, v bool) {
	rw, _ := ecs.GetResource[RealWindow](w)
	rw.RunnableOnUnfocused = v
}

func SetResizingMode(w *ecs.World, mode ResizingMode) {
	rw, _ := ecs.GetResource[RealWindow](w)
	rw.ResizingMode = mode
}

func SetWindowSizeLimits(w *ecs.World, minW, minH, maxW, maxH int) {
	rw, _ := ecs.GetResource[RealWindow](w)
	rw.SizeLimits = SizeLimits{MinW: minW, MinH: minH, MaxW: maxW, MaxH: maxH}
}

func MaximizeWindow(w *ecs.World) {
	rw, _ := ecs.GetResource[RealWindow](w)
	rw.Action = ActionMaximize
}

func MinimizeWindow(w *ecs.World) {
	rw, _ := ecs.GetResource[RealWindow](w)
	rw.Action = ActionMinimize
}

func RestoreWindow(w *ecs.World) {
	rw, _ := ecs.GetResource[RealWindow](w)
	rw.Action = ActionRestore
}

func GetRealWindowObservedSize[T constraints.Number](w *ecs.World) (width, height T) {
	rw, _ := ecs.GetResource[RealWindowObserved](w)
	width = T(rw.Width)
	height = T(rw.Height)

	return
}

func GetDeviceScale(w *ecs.World) float64 {
	rw, _ := ecs.GetResource[RealWindowObserved](w)
	return rw.DeviceScale
}

// virtual
func SetVirtualWindowSize(w *ecs.World, width, height int) {
	vw, _ := ecs.GetResource[VirtualWindow](w)
	vw.Width = width
	vw.Height = height
}

func GetVirtualWindowSize[T constraints.Number](w *ecs.World) (width, height T) {
	vw, _ := ecs.GetResource[VirtualWindow](w)
	width = T(vw.Width)
	height = T(vw.Height)

	return
}

// other
func ScreenToVirtual(w *ecs.World, sx, sy float64) (vx, vy float64, ok bool) {
	vp, _ := ecs.GetResource[Viewport](w)

	x := sx - vp.OffsetX
	y := sy - vp.OffsetY

	if vp.Scale > 0 {
		if x < 0 || y < 0 {
			return 0, 0, false
		}
		vx = x / float64(vp.Scale)
		vy = y / float64(vp.Scale)
		return vx, vy, true
	}

	if vp.ScaleF > 0 {
		if x < 0 || y < 0 {
			return 0, 0, false
		}
		vx = x / vp.ScaleF
		vy = y / vp.ScaleF
		return vx, vy, true
	}

	return 0, 0, false
}

func VirtualToScreen(w *ecs.World, vx, vy float64) (sx, sy float64) {
	vp, _ := ecs.GetResource[Viewport](w)

	if vp.Scale > 0 {
		return vx*float64(vp.Scale) + vp.OffsetX,
			vy*float64(vp.Scale) + vp.OffsetY
	}

	if vp.ScaleF > 0 {
		return vx*vp.ScaleF + vp.OffsetX,
			vy*vp.ScaleF + vp.OffsetY
	}

	return vx, vy
}
