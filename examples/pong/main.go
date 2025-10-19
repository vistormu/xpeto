package main

import (
	"pong/game"

	"github.com/vistormu/xpeto"
	"github.com/vistormu/xpeto/pkg/physics"
	"github.com/vistormu/xpeto/pkg/physics/debug"
)

func main() {
	xp.NewApp().
		WithPkgs(xp.DefaultPkgs()...).
		WithPkgs(game.Pkg).
		WithPkgs(physics.Pkg, debug.Pkg).
		Run()
}
