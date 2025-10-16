package app

type runner interface {
	run(a *App) error
}

type Runner uint8

const (
	Ebiten Runner = iota
	Headless
)
