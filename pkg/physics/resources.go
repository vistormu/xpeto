package physics

type Settings struct {
	GravityX float64
	GravityY float64
	CellSize float64
}

type DebugSettings struct {
	Enabled      bool
	ShowAABBs    bool
	ShowContacts bool
	ShowVelocity bool
	ShowGrid     bool // draw broadphase grid
}
