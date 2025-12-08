package main

import (
	"pong/game"

	"github.com/vistormu/xpeto"
	"github.com/vistormu/xpeto/backends/ebiten"
)

func main() {
	xp.NewApp(ebiten.Backend).
		AddPkg(xp.DefaultPkgs, ebiten.DefaultPkgs, game.Pkg).
		Run()
}
