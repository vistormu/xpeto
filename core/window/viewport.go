package window

import (
	"github.com/vistormu/xpeto/core/ecs"
)

type Viewport struct {
	Scale   int
	ScaleF  float64
	OffsetX float64
	OffsetY float64
}

// =======
// helpers
// =======
func getIntegerViewport(w *ecs.World) Viewport {
	obs, _ := ecs.GetResource[RealWindowObserved](w)
	vw, vh := GetVirtualWindowSize[int](w)

	realW, realH := obs.Width, obs.Height
	if vw <= 0 || vh <= 0 || realW <= 0 || realH <= 0 {
		return Viewport{Scale: 1, ScaleF: 1}
	}

	sx := realW / vw
	sy := realH / vh
	s := max(min(sx, sy), 1)

	drawW := vw * s
	drawH := vh * s

	offX := float64(realW-drawW) * 0.5
	offY := float64(realH-drawH) * 0.5

	return Viewport{
		Scale:   s,
		ScaleF:  float64(s),
		OffsetX: offX,
		OffsetY: offY,
	}
}

func getFreeViewport(w *ecs.World) Viewport {
	obs, _ := ecs.GetResource[RealWindowObserved](w)
	vw, vh := GetVirtualWindowSize[int](w)

	realW, realH := obs.Width, obs.Height
	if vw <= 0 || vh <= 0 || realW <= 0 || realH <= 0 {
		return Viewport{Scale: 1, ScaleF: 1}
	}

	sx := float64(realW) / float64(vw)
	sy := float64(realH) / float64(vh)
	s := min(sx, sy)

	drawW := float64(vw) * s
	drawH := float64(vh) * s

	offX := (float64(realW) - drawW) * 0.5
	offY := (float64(realH) - drawH) * 0.5

	return Viewport{
		Scale:   0,
		ScaleF:  s,
		OffsetX: offX,
		OffsetY: offY,
	}
}

// ===
// API
// ===
func ComputeViewport(w *ecs.World) Viewport {
	sc, _ := ecs.GetResource[Scaling](w)

	switch sc.Mode {
	case ScalingInteger:
		return getIntegerViewport(w)
	case ScalingFree:
		return getFreeViewport(w)
	case ScalingHiDPI:
		return Viewport{Scale: 1, ScaleF: 1, OffsetX: 0, OffsetY: 0}
	default:
		return getFreeViewport(w)
	}
}
