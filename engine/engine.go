package engine

import (
	"fmt"
	"io/fs"
	"reflect"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/audio"
	"github.com/vistormu/xpeto/ecs"
	"github.com/vistormu/xpeto/event"
	"github.com/vistormu/xpeto/font"
	"github.com/vistormu/xpeto/image"
	"github.com/vistormu/xpeto/input"
	"github.com/vistormu/xpeto/scene"
	"github.com/vistormu/xpeto/state"

	"github.com/vistormu/xpeto/engine/system/animation"
	audiosys "github.com/vistormu/xpeto/engine/system/audio"
	inputsys "github.com/vistormu/xpeto/engine/system/input"
	"github.com/vistormu/xpeto/engine/system/render"
	scenesys "github.com/vistormu/xpeto/engine/system/scene"
)

// ====
// game
// ====
type game struct {
	context     *ecs.Context
	accumulator float64
	lastTime    time.Time
	fixedDelta  float64
}

func newGame() *game {
	return &game{
		context:     ecs.NewContext(),
		accumulator: 0,
		lastTime:    time.Now(),
		fixedDelta:  1.0 / 60.0,
	}
}

func (g *game) Update() error {
	now := time.Now()
	frameTime := now.Sub(g.lastTime).Seconds()
	g.lastTime = now

	// Cap to avoid spiral of death
	if frameTime > 0.25 {
		frameTime = 0.25
	}

	g.accumulator += frameTime

	sm, _ := ecs.GetResource[*ecs.SystemManager](g.context)
	for _, sys := range sm.Systems() {
		if !sm.Filter(g.context, sys) {
			continue
		}

		sm.Load(g.context, sys)

		// fixed updates
		for g.accumulator >= g.fixedDelta {
			sys.FixedUpdate(g.context, float32(g.fixedDelta))
			g.accumulator -= g.fixedDelta
		}

		// variable updates
		sys.Update(g.context, float32(frameTime))
	}

	return nil
}

func (g *game) Draw(screen *image.Image) {
	sm, _ := ecs.GetResource[*ecs.SystemManager](g.context)

	for _, sys := range sm.Systems() {
		sys.Draw(screen)
	}
}

func (g *game) Layout(w, h int) (int, int) {
	return w, h
}

// ======
// engine
// ======
type Engine struct {
	game *game
}

func NewEngine() *Engine {
	e := &Engine{
		game: newGame(),
	}

	// core resources
	ecs.AddResource(e.game.context, ecs.NewEntityManager())
	ecs.AddResource(e.game.context, ecs.NewSystemManager())
	ecs.AddResource(e.game.context, event.NewManager())
	ecs.AddResource(e.game.context, input.NewManager())
	ecs.AddResource(e.game.context, scene.NewManager())
	ecs.AddResource(e.game.context, state.NewManager())

	// core systems
	var alwaysTrueFilter ecs.SystemFilter = func(*ecs.Context) bool { return true }
	e.WithSystems(
		inputsys.NewSystem(), alwaysTrueFilter,
		audiosys.NewSystem(), alwaysTrueFilter,
		scenesys.NewSystem(), alwaysTrueFilter,
		animation.NewSystem(), alwaysTrueFilter,
		render.NewSystem(), alwaysTrueFilter,
		// input
	)

	return e
}

func (e *Engine) WithAssets(fsys fs.FS) *Engine {
	ecs.AddResource(e.game.context, audio.NewManager(fsys))
	ecs.AddResource(e.game.context, image.NewManager(fsys))
	ecs.AddResource(e.game.context, font.NewManager(fsys))

	return e
}

func (e *Engine) WithScenes(scenes ...scene.Scene) *Engine {
	sm, _ := ecs.GetResource[*scene.Manager](e.game.context)

	for _, s := range scenes {
		sm.Register(s)
	}

	return e
}

func (e *Engine) WithInitialScene(s scene.Scene) *Engine {
	sm, _ := ecs.GetResource[*scene.Manager](e.game.context)

	if s == nil {
		return e
	}

	sm.Push(reflect.TypeOf(s))
	s.OnLoad(e.game.context)
	s.OnEnter(e.game.context)

	return e
}

func (e *Engine) WithSystems(systemsAndFilters ...any) *Engine {
	sm, _ := ecs.GetResource[*ecs.SystemManager](e.game.context)

	for i := 0; i < len(systemsAndFilters); i += 2 {
		system, ok := systemsAndFilters[i].(ecs.System)
		if !ok {
			panic(fmt.Sprintf("expected a System type, got %T", systemsAndFilters[i]))
		}
		filter, ok := systemsAndFilters[i+1].(ecs.SystemFilter)
		if !ok {
			panic(fmt.Sprintf("expected a SystemFilter type, got %T", systemsAndFilters[i+1]))
		}

		sm.Register(system, filter)
	}

	return e
}

// tmp
func (e *Engine) WithKeys(keys ...input.Key) *Engine {
	im, _ := ecs.GetResource[*input.Manager](e.game.context)
	im.Register(keys...)

	return e
}

func (e *Engine) Run() {
	ebiten.SetWindowSize(800, 600)
	err := ebiten.RunGame(e.game)
	if err != nil {
		panic(err)
	}
}
