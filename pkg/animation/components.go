package animation

import (
	"github.com/vistormu/xpeto/pkg/render"
)

type Animation struct {
	Frames   []render.Image
	Duration float32
	Elapsed  float32
	Current  uint64
	Loop     bool
}
