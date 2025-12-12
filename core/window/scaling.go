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

// ===
// API
// ===
func SetScalingMode(w *ecs.World, mode ScalingMode) {
	s, _ := ecs.GetResource[Scaling](w)
	s.Mode = mode
}

func SetPixelSnap(w *ecs.World, v bool) {
	s, _ := ecs.GetResource[Scaling](w)
	s.SnapPixels = v
}

func GetDesiredVirtualSize(w *ecs.World) (vw, vh int, ok bool) {
	sc, _ := ecs.GetResource[Scaling](w)
	if sc.Mode != ScalingHiDPI {
		return 0, 0, false
	}

	obs, _ := ecs.GetResource[RealWindowObserved](w)
	vw = int(float64(obs.Width)*obs.DeviceScale + 0.5)
	vh = int(float64(obs.Height)*obs.DeviceScale + 0.5)

	return vw, vh, true
}
