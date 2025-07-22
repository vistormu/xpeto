package engine

import (
	"fmt"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/internal/ecs"
	"github.com/vistormu/xpeto/internal/event"
	"github.com/vistormu/xpeto/internal/state"
)

// ====
// game
// ====
type ebitenGame struct {
	context      *ecs.Context
	stateEngines []state.StateSystem
}

func newGame() *ebitenGame {
	return &ebitenGame{
		context:      ecs.NewContext(),
		stateEngines: make([]state.StateSystem, 0),
	}
}

func (g *ebitenGame) Update() error {
	for _, engine := range g.stateEngines {
		engine.Update(g.context)
	}

	return nil
}

func (g *ebitenGame) Draw(screen *ebiten.Image) {
}

func (g *ebitenGame) Layout(w, h int) (int, int) {
	// sm, _ := ecs.GetResource[*scene.Manager](g.context)

	// current, ok := sm.Current()
	// if !ok {
	// 	return w, h
	// }

	// return current.Layout(w, h)
	return 800, 600
}

// =======
// builder
// =======
type Game struct {
	fsys fs.FS

	resources []any
	systems   []struct {
		hook  state.Hook
		state any
		fn    any
	}
}

func NewGame() *Game {
	return &Game{
		fsys:      nil,
		resources: []any{},
		systems: []struct {
			hook  state.Hook
			state any
			fn    any
		}{},
	}
}

func (g *Game) WithResources(resources ...any) *Game {
	g.resources = append(g.resources, resources...)
	return g
}

func (g *Game) WithAssets(fsys fs.FS) *Game {
	g.fsys = fsys
	return g
}

func (g *Game) WithSystems(hook state.Hook, st any, fn any) *Game {
	g.systems = append(g.systems, struct {
		hook  state.Hook
		state any
		fn    any
	}{
		hook:  hook,
		state: st,
		fn:    fn,
	})

	return g
}

func (g *Game) WithPkg(pkg Pkg) *Game {
	return g
}

func (g *Game) Run() {
	game := newGame()

	// core resources
	ecs.AddResource(game.context, ecs.NewManager())
	ecs.AddResource(game.context, event.NewManager())

	// user resources
	for _, res := range g.resources {
		ecs.AddResource(game.context, res)
	}

	// systems
	registered := make(map[any]state.StateManager)
	for _, s := range g.systems {
		sm, ok := registered[s.state]
		if !ok {
			sm = state.NewManager(s.state)
		}

		sm.Register(s.hook, s.state, s.fn)
	}

	for st, sm := range registered {
		ecs.AddResource(game.context, sm)
		ss := state.NewSystem(st)
		ss.OnEnter(game.context)
		game.stateEngines = append(game.stateEngines, ss)
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Xpeto Game Engine")

	err := ebiten.RunGame(game)
	if err != nil {
		panic(fmt.Sprintf("failed to run game: %v", err))
	}
}
