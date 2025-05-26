package audio

import (
	"github.com/vistormu/xpeto/audio"
)

type AudioPlay struct {
	Audio  audio.Handle
	Loop   bool
	Volume float32
}
