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
	stateEngine  *state.System
	renderEngine *render.System
}

func newGame() *ebitenGame {
	return &ebitenGame{
		context:      ecs.NewContext(),
		stateEngine:  nil,
		renderEngine: nil,
	}
}

func (g *ebitenGame) Update() error {
	g.stateEngine.Update(g.context)
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
	systems   []struct {
		hook  state.Hook
		state comparable
		fn    any
	}
}

func NewGame() *Game {
	return &Game{
		fsys:      nil,
		resources: []any{},
		systems: []struct {
			state  comparable
			fnType state.StateFn
			fn     any
		}{},
	}
}

func (gb *Game) WithResources(resources ...any) *Game {
	gb.resources = append(gb.resources, resources...)
	return gb
}

func (gb *Game) WithAssets(fsys fs.FS) *Game {
	gb.fsys = fsys
	return gb
}

func (gb *Game) WithSystems(fnType state.StateFn, st state.State, fn any) *Game {
	gb.functions = append(gb.functions, struct {
		state  state.State
		fnType state.StateFn
		fn     any
	}{
		state:  st,
		fnType: fnType,
		fn:     fn,
	})
	return gb
}

func (gb *Game) WithPlugin(plugin Plugin) *Game {
	return gb
}

func (gb *Game) Run() {
	game := newGame()

	// resources
	coreResources := []any{
		ecs.NewManager(),
		event.NewManager(),
		state.NewManager(),

		// tmp
		render.NewManager().WithFilesystem(gb.fsys),
	}

	for _, res := range coreResources {
		ecs.AddResource(game.context, res)
	}

	for _, res := range gb.resources {
		ecs.AddResource(game.context, res)
	}

	// core systems
	game.stateEngine = state.NewSystem()
	game.renderEngine = render.NewSystem()

	// states
	for _, st := range gb.states {
	}

	// state functions
	for _, fn := range gb.functions {
		ecs.MustResource[*state.Manager](game.context).RegisterFn(fn.state, fn.fn, fn.fnType)
	}
	game.stateEngine.OnEnter(game.context)

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Xpeto Game Engine")

	err := ebiten.RunGame(game)
	if err != nil {
		panic(fmt.Sprintf("failed to run game: %v", err))
	}
}
