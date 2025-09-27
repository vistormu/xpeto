package xp

import (
	"reflect"

	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"
	"github.com/vistormu/xpeto/internal/engine"
	"github.com/vistormu/xpeto/internal/event"
	"github.com/vistormu/xpeto/internal/scheduler"

	"github.com/vistormu/xpeto/pkg"
	"github.com/vistormu/xpeto/pkg/asset"
	"github.com/vistormu/xpeto/pkg/image"
	"github.com/vistormu/xpeto/pkg/text"
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

func AddResourceByType(ctx *Context, resource any, type_ reflect.Type) {
	core.AddResourceByType(ctx, resource, type_)
}

func GetResource[T any](ctx *Context) (T, bool) {
	return core.GetResource[T](ctx)
}

func MustResource[T any](ctx *Context) T {
	return core.MustResource[T](ctx)
}

// plugin
type Plugin = core.Plugin
type ScheduleBuilder = core.ScheduleBuilder

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

func And(filters ...Filter) Filter {
	return ecs.And(filters...)
}

func Or(filters ...Filter) Filter {
	return ecs.Or(filters...)
}

func Not(filter Filter) Filter {
	return ecs.Not(filter)
}

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

// ======
// engine
// ======
// settings
type GameSettings = engine.Settings

// engine
type Game = engine.Game

func NewGame() *Game {
	return engine.NewGame()
}

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
type Stage = core.Stage
type Schedule = scheduler.Schedule

const (
	PreStartup  = core.PreStartup
	Startup     = core.Startup
	PostStartup = core.PostStartup

	First     = core.First
	PreUpdate = core.PreUpdate

	FixedFirst      = core.FixedFirst
	FixedPreUpdate  = core.FixedPreUpdate
	FixedUpdate     = core.FixedUpdate
	FixedPostUpdate = core.FixedPostUpdate
	FixedLast       = core.FixedLast

	Update     = core.Update
	PostUpdate = core.PostUpdate
	Last       = core.Last
)

// scheduler
// type Scheduler = scheduler.Scheduler

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
// fonts
// =====
// types
type Font = text.Font

// components
type Text = text.Text

// =========
// transform
// =========

// components
type Transform = transform.Transform
