package engine

import (
	"fmt"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/internal/ecs"
	"github.com/vistormu/xpeto/internal/event"
	"github.com/vistormu/xpeto/internal/state"
	"github.com/vistormu/xpeto/pkg/render"
)

// ====
// game
// ====
type ebitenGame struct {
	context      *ecs.Context
	stateEngines []state.StateSystem
	renderEngine *render.System
}

func newGame() *ebitenGame {
	return &ebitenGame{
		context:      ecs.NewContext(),
		stateEngines: make([]state.StateSystem, 0),
		renderEngine: nil,
	}
}

func (g *ebitenGame) Update() error {
	for _, engine := range g.stateEngines {
		engine.Update(g.context)
	}

	g.renderEngine.Update(g.context)

	return nil
}

func (g *ebitenGame) Draw(screen *ebiten.Image) {
	g.renderEngine.Draw(screen)
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
	states    []any
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
		states:    []any{},
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

func (g *Game) WithState(state any) *Game {
	g.states = append(g.states, state)
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

	// states
	for _, st := range g.states {
		ecs.AddResource(game.context, state.NewManager(st))
		// game.stateEngines = append(game.stateEngines, state.NewSystem[]())
	}

	// render
	// game.stateEngine = state.NewSystem()
	// game.renderEngine = render.NewSystem()
	//
	// // state functions
	// for _, fn := range g.functions {
	// 	ecs.MustResource[*state.Manager](game.context).RegisterFn(fn.state, fn.fn, fn.fnType)
	// }
	// game.stateEngine.OnEnter(game.context)

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Xpeto Game Engine")

	err := ebiten.RunGame(game)
	if err != nil {
		panic(fmt.Sprintf("failed to run game: %v", err))
	}
}
