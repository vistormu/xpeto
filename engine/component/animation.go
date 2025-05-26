package component

import (
	"github.com/vistormu/xpeto/image"
)

type Animation struct {
	Frames   []image.Handle
	Duration float32
	Elapsed  float32
	Current  uint64
	Loop     bool
}
