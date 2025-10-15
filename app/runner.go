package app

type runner interface {
	run(a *App) error
}

type Runner uint8

const (
	Ebiten Runner = iota
	Headless
)

var toRunner = map[Runner]runner{
	Ebiten:   &ebitenRunner{},
	Headless: &headlessRunner{},
}
