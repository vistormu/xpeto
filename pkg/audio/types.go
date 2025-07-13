package audio

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/vistormu/xpeto/internal/core"
)

type Player = audio.Player
type Context = audio.Context
type Audio = core.Handle[*Player]
