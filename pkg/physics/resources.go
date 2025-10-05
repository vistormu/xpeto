package physics

import (
	"github.com/vistormu/xpeto/internal/core"
)

type Settings struct {
	Gravity  core.Vector[float32]
	CellSize float64
}

type DebugSettings struct {
	Enabled      bool
	ShowAABBs    bool
	ShowContacts bool
	ShowVelocity bool
	ShowGrid     bool // draw broadphase grid
}
