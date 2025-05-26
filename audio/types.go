package audio

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/vistormu/xpeto/internal/core"
)

type Audio = audio.Player
type Context = audio.Context
type Handle = core.Handle[*Audio]
