package animation

import (
	"github.com/vistormu/xpeto/pkg/asset"
)

type Animation struct {
	Frames   []asset.Handle
	Duration float32
	Elapsed  float32
	Current  uint64
	Loop     bool
}
