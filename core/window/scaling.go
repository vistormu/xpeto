package window

import (
	"github.com/vistormu/xpeto/core/ecs"
)

type ScalingMode uint8

const (
	ScalingFree ScalingMode = iota
	ScalingInteger
	ScalingHiDPI
)

type Scaling struct {
	Mode       ScalingMode
	SnapPixels bool
}

func newScaling() Scaling {
	return Scaling{
		Mode:       ScalingInteger,
		SnapPixels: true,
	}
}

// ===
// API
// ===
func SetScalingMode(w *ecs.World, mode ScalingMode) {
	s, ok := ecs.GetResource[Scaling](w)
	if !ok {
		return
	}
	s.Mode = mode
}

func SetPixelSnap(w *ecs.World, v bool) {
	s, ok := ecs.GetResource[Scaling](w)
	if !ok {
		return
	}
	s.SnapPixels = v
}

func GetDesiredVirtualSize(w *ecs.World) (vw, vh int, ok bool) {
	s, ok := ecs.GetResource[Scaling](w)
	if !ok {
		return
	}
	if s.Mode != ScalingHiDPI {
		return 0, 0, false
	}

	obs, ok := ecs.GetResource[RealWindowObserved](w)
	if !ok {
		return
	}

	vw = int(float64(obs.Width)*obs.DeviceScale + 0.5)
	vh = int(float64(obs.Height)*obs.DeviceScale + 0.5)

	return vw, vh, true
}
