package window

import (
	"testing"

	"github.com/vistormu/xpeto/core/ecs"
)

func newWorldForWindow() *ecs.World {
	w := ecs.NewWorld()

	ecs.AddResource(w, RealWindowObserved{Width: 800, Height: 600, DeviceScale: 1})
	ecs.AddResource(w, VirtualWindow{Width: 400, Height: 300})
	ecs.AddResource(w, Scaling{Mode: ScalingInteger, SnapPixels: true})

	// mimic plugin insertion
	ecs.AddResource(w, ComputeViewport(w))

	return w
}

func TestComputeViewport_IntegerScaling_Normal(t *testing.T) {
	w := newWorldForWindow()

	vp := ComputeViewport(w)
	if vp.Scale != 2 {
		t.Fatalf("expected integer scale 2, got %d (ScaleF=%v)", vp.Scale, vp.ScaleF)
	}
	if vp.OffsetX != 0 || vp.OffsetY != 0 {
		t.Fatalf("expected offsets 0,0 got %v,%v", vp.OffsetX, vp.OffsetY)
	}
}

func TestComputeViewport_FreeScaling_Normal(t *testing.T) {
	w := newWorldForWindow()
	SetScalingMode(w, ScalingFree)

	vp := ComputeViewport(w)
	if vp.Scale != 0 {
		t.Fatalf("expected Scale=0 in free scaling, got %d", vp.Scale)
	}
	if vp.ScaleF <= 0 {
		t.Fatalf("expected ScaleF>0 in free scaling, got %v", vp.ScaleF)
	}
	if vp.OffsetX != 0 || vp.OffsetY != 0 {
		t.Fatalf("expected offsets 0,0 got %v,%v", vp.OffsetX, vp.OffsetY)
	}
}

func TestComputeViewport_InvalidSizes_ReturnsSafeDefaults(t *testing.T) {
	w := ecs.NewWorld()

	ecs.AddResource(w, RealWindowObserved{Width: 0, Height: 0, DeviceScale: 1})
	ecs.AddResource(w, VirtualWindow{Width: 0, Height: 0})
	ecs.AddResource(w, Scaling{Mode: ScalingInteger, SnapPixels: true})

	vp := ComputeViewport(w)
	if vp.Scale != 1 {
		t.Fatalf("expected safe default Scale=1, got %d", vp.Scale)
	}
	if vp.ScaleF != 1 {
		t.Fatalf("expected safe default ScaleF=1, got %v", vp.ScaleF)
	}
	if vp.OffsetX != 0 || vp.OffsetY != 0 {
		t.Fatalf("expected safe default offsets 0,0 got %v,%v", vp.OffsetX, vp.OffsetY)
	}
}

func TestComputeViewport_IntegerScaling_RealSmallerThanVirtual_FallsBackToFree(t *testing.T) {
	w := ecs.NewWorld()

	ecs.AddResource(w, RealWindowObserved{Width: 320, Height: 240, DeviceScale: 1})
	ecs.AddResource(w, VirtualWindow{Width: 640, Height: 480})
	ecs.AddResource(w, Scaling{Mode: ScalingInteger, SnapPixels: true})

	vp := ComputeViewport(w)

	// current behaviour: fallback to free scaling
	if vp.Scale != 0 {
		t.Fatalf("expected fallback to free scaling (Scale=0), got Scale=%d ScaleF=%v", vp.Scale, vp.ScaleF)
	}
	if vp.ScaleF <= 0 {
		t.Fatalf("expected ScaleF>0 on fallback, got %v", vp.ScaleF)
	}
	if vp.OffsetX < 0 || vp.OffsetY < 0 {
		t.Fatalf("expected non-negative offsets, got %v,%v", vp.OffsetX, vp.OffsetY)
	}
}

func TestScreenToVirtual_RejectsLetterboxArea_IntegerScaling(t *testing.T) {
	w := ecs.NewWorld()

	// real is wider, so we get horizontal letterboxing with integer scaling
	ecs.AddResource(w, RealWindowObserved{Width: 1000, Height: 600, DeviceScale: 1})
	ecs.AddResource(w, VirtualWindow{Width: 400, Height: 300})
	ecs.AddResource(w, Scaling{Mode: ScalingInteger, SnapPixels: true})

	vp := ComputeViewport(w)
	ecs.AddResource(w, vp)

	// integer scale is 2 -> draw size 800x600, OffsetX should be (1000-800)/2 = 100
	if vp.Scale != 2 {
		t.Fatalf("expected scale=2, got %d", vp.Scale)
	}
	if vp.OffsetX != 100 {
		t.Fatalf("expected OffsetX=100, got %v", vp.OffsetX)
	}

	// point in left letterbox should be rejected
	_, _, ok := ScreenToVirtual(w, 50, 10)
	if ok {
		t.Fatalf("expected ok=false for letterbox area")
	}

	// point in right letterbox should be rejected
	_, _, ok = ScreenToVirtual(w, 950, 10)
	if ok {
		t.Fatalf("expected ok=false for letterbox area")
	}

	// point inside draw area should be accepted
	vx, vy, ok := ScreenToVirtual(w, 110, 20)
	if !ok {
		t.Fatalf("expected ok=true inside draw area")
	}
	if vx < 0 || vy < 0 {
		t.Fatalf("expected non-negative virtual coords, got %v,%v", vx, vy)
	}
}

func TestVirtualToScreen_ScreenToVirtual_RoundTrip_IntegerScaling(t *testing.T) {
	w := ecs.NewWorld()

	ecs.AddResource(w, RealWindowObserved{Width: 800, Height: 600, DeviceScale: 1})
	ecs.AddResource(w, VirtualWindow{Width: 400, Height: 300})
	ecs.AddResource(w, Scaling{Mode: ScalingInteger, SnapPixels: true})

	vp := ComputeViewport(w)
	ecs.AddResource(w, vp)

	sx, sy := VirtualToScreen(w, 10, 20)
	vx, vy, ok := ScreenToVirtual(w, sx, sy)
	if !ok {
		t.Fatalf("expected round-trip ok=true")
	}

	// exact in integer scaling
	if vx != 10 || vy != 20 {
		t.Fatalf("expected round-trip (10,20), got (%v,%v)", vx, vy)
	}
}

func TestSetWindowSizeLimits_Normalises(t *testing.T) {
	w := ecs.NewWorld()
	ecs.AddResource(w, RealWindow{})

	SetWindowSizeLimits(w, -10, -20, 0, -1)

	rw, _ := ecs.GetResource[RealWindow](w)
	if rw.SizeLimits.MinW != -1 || rw.SizeLimits.MinH != -1 {
		t.Fatalf("expected negative mins normalised to -1, got %+v", rw.SizeLimits)
	}
	if rw.SizeLimits.MaxW != -1 || rw.SizeLimits.MaxH != -1 {
		t.Fatalf("expected non-positive max normalised to -1, got %+v", rw.SizeLimits)
	}

	// inverted bounds swap
	SetWindowSizeLimits(w, 800, 600, 200, 100)
	rw, _ = ecs.GetResource[RealWindow](w)
	if rw.SizeLimits.MinW != 200 || rw.SizeLimits.MaxW != 800 {
		t.Fatalf("expected swapped W bounds, got %+v", rw.SizeLimits)
	}
	if rw.SizeLimits.MinH != 100 || rw.SizeLimits.MaxH != 600 {
		t.Fatalf("expected swapped H bounds, got %+v", rw.SizeLimits)
	}
}

func TestSetRealWindowSize_ClampsToAtLeastOne(t *testing.T) {
	w := ecs.NewWorld()
	ecs.AddResource(w, RealWindow{})

	SetRealWindowSize(w, 0, -5)

	rw, _ := ecs.GetResource[RealWindow](w)
	if rw.Width != 1 || rw.Height != 1 {
		t.Fatalf("expected clamped size 1x1, got %dx%d", rw.Width, rw.Height)
	}
}

func TestSetVirtualWindowSize_ClampsToAtLeastOne(t *testing.T) {
	w := ecs.NewWorld()
	ecs.AddResource(w, VirtualWindow{})

	SetVirtualWindowSize(w, 0, -5)

	vw, _ := ecs.GetResource[VirtualWindow](w)
	if vw.Width != 1 || vw.Height != 1 {
		t.Fatalf("expected clamped size 1x1, got %dx%d", vw.Width, vw.Height)
	}
}
