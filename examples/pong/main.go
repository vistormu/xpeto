package main

import (
	"pong/game"

	"github.com/vistormu/xpeto"
)

func main() {
	xp.NewApp().
		AddPkg(xp.DefaultPkgs).
		AddPkg(game.Pkg).
		Run()
}
