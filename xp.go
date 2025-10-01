package xp

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"
	"github.com/vistormu/xpeto/internal/event"
	"github.com/vistormu/xpeto/internal/game"
	"github.com/vistormu/xpeto/internal/schedule"

	"github.com/vistormu/xpeto/pkg"
	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/graphics"
	"github.com/vistormu/xpeto/pkg/image"
	"github.com/vistormu/xpeto/pkg/input"
	"github.com/vistormu/xpeto/pkg/render"
	"github.com/vistormu/xpeto/pkg/text"
	"github.com/vistormu/xpeto/pkg/time"
	"github.com/vistormu/xpeto/pkg/transform"
)

// ============================================
// these re-exports are part of the core engine
// ============================================

// ====
// core
// ====
// context
type Context = core.Context

func AddResource[T any](ctx *Context, resource T) {
	core.AddResource(ctx, resource)
}

var AddResourceByType = core.AddResourceByType

func GetResource[T any](ctx *Context) (T, bool) {
	return core.GetResource[T](ctx)
}

func MustResource[T any](ctx *Context) T {
	return core.MustResource[T](ctx)
}

// geometry
type Vector = core.Vector[float32]

// ===
// ecs
// ===
// types
type Entity = ecs.Entity
type Component = ecs.Component
type System = ecs.System

// filters
type Filter = ecs.Filter

func Has[T any]() Filter {
	return ecs.Has[T]()
}

var And = ecs.And
var Or = ecs.Or
var Not = ecs.Not

// world
type World = ecs.World

func CreateEntity(ctx *Context) Entity {
	w := MustResource[*World](ctx)
	return w.Create()
}

func DestroyEntity(ctx *Context, entity Entity) {
	w := MustResource[*World](ctx)
	w.Destroy(entity)
}

func DestroyAllEntities(ctx *Context) {
	w := MustResource[*World](ctx)
	w.DestroyAll()
}

func Query(ctx *Context, f Filter) []Entity {
	w := MustResource[*World](ctx)
	return w.Query(f)
}

func GetComponent[T any](ctx *Context, entity Entity) (T, bool) {
	w := MustResource[*World](ctx)
	return ecs.GetComponent[T](w, entity)
}

func AddComponent[T any](ctx *Context, entity Entity, component T) {
	w := MustResource[*World](ctx)
	ecs.AddComponent(w, entity, component)
}

func RemoveComponent[T any](ctx *Context, entity Entity) {
	w := MustResource[*World](ctx)
	ecs.RemoveComponent[T](w, entity)
}

// ====
// game
// ====
// settings
type GameSettings = game.Settings

// engine
type Game = game.Game

var NewGame = game.NewGame

// plugin
type Plugin game.Plugin

// =====
// event
// =====
// types
type Event = event.Event

// bus
type EventBus = event.Bus

func Subscribe[T any](ctx *Context, callback func(data T)) Event {
	eb := core.MustResource[*event.Bus](ctx)
	return event.Subscribe(eb, callback)
}

func Unsubscribe(ctx *Context, e Event) {
	eb := core.MustResource[*event.Bus](ctx)
	eb.Unsubscribe(e)
}

func Publish[T any](ctx *Context, data T) {
	eb := core.MustResource[*event.Bus](ctx)
	event.Publish(eb, data)
}

// =========
// scheduler
// =========
// types
type Stage = schedule.Stage
type Schedule = schedule.Schedule
type Scheduler = schedule.Scheduler
type State[T comparable] = schedule.State[T]
type NextState[T comparable] = schedule.NextState[T]

var PreStartup = schedule.PreStartup
var Startup = schedule.Startup
var PostStartup = schedule.PostStartup

var First = schedule.First
var PreUpdate = schedule.PreUpdate

func OnExit[T comparable](state T) Stage {
	return schedule.OnExit(state)
}
func OnTransition[T comparable](from, to T) Stage {
	return schedule.OnTransition(from, to)
}
func OnEnter[T comparable](state T) Stage {
	return schedule.OnEnter(state)
}

var FixedFirst = schedule.FixedFirst
var FixedPreUpdate = schedule.FixedPreUpdate
var FixedUpdate = schedule.FixedUpdate
var FixedPostUpdate = schedule.FixedPostUpdate
var FixedLast = schedule.FixedLast

var Update = schedule.Update
var PostUpdate = schedule.PostUpdate
var Last = schedule.Last

var PreDraw = schedule.PreDraw
var Draw = schedule.Draw
var PostDraw = schedule.PostDraw

// conditions
func InState[T comparable](s T) func(*core.Context) bool {
	return schedule.InState(s)
}

// scheduler
func AddStateMachine[T comparable](sch *Scheduler, initial T) {
	schedule.AddStateMachine(sch, initial)
}

var AddSystem = schedule.AddSystem

// states
func CurrentState[T comparable](ctx *Context) T {
	current, ok := core.GetResource[*State[T]](ctx)
	if !ok {
		var zero T
		return zero
	}

	return current.Get()
}

func SetNextState[T comparable](ctx *Context, s T) {
	next, ok := core.GetResource[*NextState[T]](ctx)
	if !ok {
		return
	}

	next.Set(s)
}

// events
type EventStateTransition[T comparable] = schedule.EventStateTransition[T]

// ============================================
// these re-exports are part of opt-in packages
// ============================================
var DefaultPlugins = pkg.DefaultPlugins

// =====
// asset
// =====
// types
type Handle = asset.Handle
type LoadState = asset.LoadState
type LoaderFn = asset.LoaderFn

const (
	NotFound = asset.NotFound
	Loading  = asset.Loading
	Loaded   = asset.Loaded
	Failed   = asset.Failed
)

// server
type AssetServer = asset.Server

func AddAssetLoader(ctx *Context, ext string, loader LoaderFn) {
	as, ok := core.GetResource[*asset.Server](ctx)
	if !ok {
		return
	}
	as.AddLoader(ext, loader)
}

func LoadAsset[T any, B any](ctx *Context) {
	as, ok := core.GetResource[*AssetServer](ctx)
	if !ok {
		return
	}
	asset.Load[T, B](as)
}

func GetAsset[T any](ctx *Context, handle Handle) (T, bool) {
	as, ok := core.GetResource[*AssetServer](ctx)
	if !ok {
		var zero T
		return zero, false
	}
	return asset.GetAsset[T](as, handle)
}

func GetState(ctx *Context, handle Handle) LoadState {
	as, ok := core.GetResource[*AssetServer](ctx)
	if !ok {
		return NotFound
	}
	return as.GetState(handle)
}

// event
type AssetEvent = asset.AssetEvent
type AssetEventKind = asset.AssetEventKind

const (
	AssetAdded    = asset.Added
	AssetModified = asset.Modified
	AssetRemoved  = asset.Removed
)

// plugin
var AssetPlugin = asset.AssetPlugin

// =====
// image
// =====
// types
type Image = image.Image

// components
type Sprite = image.Sprite

// =====
// input
// =====
// types
type Key = input.Key
type Keyboard = input.Keyboard
type MouseButton = input.MouseButton
type Mouse = input.Mouse
type GamepadButton = input.GamepadButton
type GamepadAxis = input.GamepadAxis
type Gamepad = input.Gamepad

const (
	KeyA Key = ebiten.KeyA
	KeyB Key = ebiten.KeyB

	KeyEnter Key = ebiten.KeyEnter
)

// events
type KeyJustPressed = input.KeyJustPressed
type KeyJustReleased = input.KeyJustReleased
type MouseButtonJustPressed = input.MouseButtonJustPressed
type MouseButtonJustReleased = input.MouseButtonJustReleased

// =====
// fonts
// =====
// types
type Font = text.Font

// components
type Text = text.Text

// ========
// graphics
// ========
type Circle = graphics.Circle

// ======
// render
// ======
// components
type Renderable = render.Renderable

// ====
// time
// ====
// types
type Time = time.Time

// =========
// transform
// =========

// components
type Transform = transform.Transform
