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

func newRealWindow() RealWindow {
	return RealWindow{
		Title:               "xpeto app",
		Width:               800,
		Height:              600,
		FullScreen:          false,
		AntiAliasing:        false,
		VSync:               false,
		RunnableOnUnfocused: true,
		ResizingMode:        ResizingModeDisabled,
		SizeLimits:          SizeLimits{-1, -1, -1, -1},
		Action:              ActionNone,
	}
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

func newVirtualWindow() VirtualWindow {
	return VirtualWindow{
		Width:  800,
		Height: 600,
	}
}

// ===
// API
// ===

// real
func SetRealWindowSize(w *ecs.World, width, height int) {
	rw, ok := ecs.GetResource[RealWindow](w)
	if !ok {
		return
	}
	rw.Width = max(width, 1)
	rw.Height = max(height, 1)
}

func GetRealWindowSize[T constraints.Number](w *ecs.World) (width, height T) {
	rw, ok := ecs.GetResource[RealWindow](w)
	if !ok {
		return
	}
	width = T(rw.Width)
	height = T(rw.Height)

	return
}

func SetFullScreen(w *ecs.World, v bool) {
	rw, ok := ecs.GetResource[RealWindow](w)
	if !ok {
		return
	}
	rw.FullScreen = v
}

func SetAntiAliasing(w *ecs.World, v bool) {
	rw, ok := ecs.GetResource[RealWindow](w)
	if !ok {
		return
	}
	rw.AntiAliasing = v
}

func SetVSync(w *ecs.World, v bool) {
	rw, ok := ecs.GetResource[RealWindow](w)
	if !ok {
		return
	}
	rw.VSync = v
}

func SetRunnableOnUnfocused(w *ecs.World, v bool) {
	rw, ok := ecs.GetResource[RealWindow](w)
	if !ok {
		return
	}
	rw.RunnableOnUnfocused = v
}

func SetResizingMode(w *ecs.World, mode ResizingMode) {
	rw, ok := ecs.GetResource[RealWindow](w)
	if !ok {
		return
	}
	rw.ResizingMode = mode
}

func SetWindowSizeLimits(w *ecs.World, minW, minH, maxW, maxH int) {
	rw, ok := ecs.GetResource[RealWindow](w)
	if !ok {
		return
	}

	if minW < 0 {
		minW = -1
	}
	if minH < 0 {
		minH = -1
	}
	if maxW <= 0 {
		maxW = -1
	}
	if maxH <= 0 {
		maxH = -1
	}
	if minW > 0 && maxW > 0 && minW > maxW {
		minW, maxW = maxW, minW
	}
	if minH > 0 && maxH > 0 && minH > maxH {
		minH, maxH = maxH, minH
	}

	rw.SizeLimits = SizeLimits{MinW: minW, MinH: minH, MaxW: maxW, MaxH: maxH}
}

func MaximizeWindow(w *ecs.World) {
	rw, ok := ecs.GetResource[RealWindow](w)
	if !ok {
		return
	}
	rw.Action = ActionMaximize
}

func MinimizeWindow(w *ecs.World) {
	rw, ok := ecs.GetResource[RealWindow](w)
	if !ok {
		return
	}
	rw.Action = ActionMinimize
}

func RestoreWindow(w *ecs.World) {
	rw, ok := ecs.GetResource[RealWindow](w)
	if !ok {
		return
	}
	rw.Action = ActionRestore
}

func GetRealWindowObservedSize[T constraints.Number](w *ecs.World) (width, height T) {
	rw, ok := ecs.GetResource[RealWindowObserved](w)
	if !ok {
		return
	}
	width = T(rw.Width)
	height = T(rw.Height)

	return
}

func GetDeviceScale(w *ecs.World) float64 {
	rw, ok := ecs.GetResource[RealWindowObserved](w)
	if !ok {
		return 0
	}
	return rw.DeviceScale
}

// virtual
func SetVirtualWindowSize(w *ecs.World, width, height int) {
	vw, ok := ecs.GetResource[VirtualWindow](w)
	if !ok {
		return
	}
	vw.Width = max(width, 1)
	vw.Height = max(height, 1)
}

func GetVirtualWindowSize[T constraints.Number](w *ecs.World) (width, height T) {
	vw, ok := ecs.GetResource[VirtualWindow](w)
	if !ok {
		return
	}
	width = T(vw.Width)
	height = T(vw.Height)

	return
}

// other
func ScreenToVirtual(w *ecs.World, sx, sy float64) (vx, vy float64, ok bool) {
	vp, ok := ecs.GetResource[Viewport](w)
	if !ok {
		return
	}

	vw, vh := GetVirtualWindowSize[int](w)
	if vw <= 0 || vh <= 0 {
		return 0, 0, false
	}

	x := sx - vp.OffsetX
	y := sy - vp.OffsetY

	if vp.Scale > 0 {
		drawW := float64(vw * vp.Scale)
		drawH := float64(vh * vp.Scale)
		if x < 0 || y < 0 || x >= drawW || y >= drawH {
			return 0, 0, false
		}
		vx = x / float64(vp.Scale)
		vy = y / float64(vp.Scale)
		return vx, vy, true
	}

	if vp.ScaleF > 0 {
		drawW := float64(vw) * vp.ScaleF
		drawH := float64(vh) * vp.ScaleF
		if x < 0 || y < 0 || x >= drawW || y >= drawH {
			return 0, 0, false
		}
		vx = x / vp.ScaleF
		vy = y / vp.ScaleF
		return vx, vy, true
	}

	return 0, 0, false
}

func VirtualToScreen(w *ecs.World, vx, vy float64) (sx, sy float64) {
	vp, ok := ecs.GetResource[Viewport](w)
	if !ok {
		return
	}

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
