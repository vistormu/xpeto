package xp

import (
	"reflect"

	"github.com/vistormu/xpeto/internal/core"
	"github.com/vistormu/xpeto/internal/ecs"
	"github.com/vistormu/xpeto/internal/engine"
	"github.com/vistormu/xpeto/internal/event"
	"github.com/vistormu/xpeto/internal/scheduler"

	"github.com/vistormu/xpeto/pkg"
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

// ===
// ecs
// ===
// types
type Entity = ecs.Entity
type Component = ecs.Component
type Archetype = ecs.Archetype
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

func CreateEntity(ctx *Context, archetype Archetype) Entity {
	w := MustResource[*World](ctx)
	return w.Create(archetype)
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
