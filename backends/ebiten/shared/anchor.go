package shared

import (
	"github.com/vistormu/xpeto/pkg/render"
)

// transforms any anchor to ebiten's top-left
func Offset(w, h float64, a render.Anchor) (dx, dy float64) {
	switch a {
	case render.AnchorCenter:
		return -w * 0.5, -h * 0.5
	case render.AnchorTopLeft:
		return 0, 0
	case render.AnchorTop:
		return -w * 0.5, 0
	case render.AnchorTopRight:
		return -w, 0
	case render.AnchorLeft:
		return 0, -h * 0.5
	case render.AnchorRight:
		return -w, -h * 0.5
	case render.AnchorBottomLeft:
		return 0, -h
	case render.AnchorBottom:
		return -w * 0.5, -h
	case render.AnchorBottomRight:
		return -w, -h
	default:
		return -w * 0.5, -h * 0.5
	}
}
