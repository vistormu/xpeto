package xp

import (
	"github.com/vistormu/xpeto/core/ecs"
	"github.com/vistormu/xpeto/core/event"
	"github.com/vistormu/xpeto/core/game"
	"github.com/vistormu/xpeto/core/schedule"
	"github.com/vistormu/xpeto/pkg/state"
)

// #############
// CORE FEATURES
// #############

// ===
// ecs
// ===

// the world is the central
type World = ecs.World

// ------
// entity
// ------

// an entity is...
// packed [index:gen]
type Entity = ecs.Entity

// this function checks
var AddEntity = ecs.AddEntity

// description
var RemoveEntity = ecs.RemoveEntity

// description
var HasEntity = ecs.HasEntity

// ---------
// component
// ---------

func AddComponent[T any](w *World, e Entity, c T) bool {
	return ecs.AddComponent(w, e, c)
}

func GetComponent[T any](w *World, e Entity) (*T, bool) {
	return ecs.GetComponent[T](w, e)
}

func RemoveComponent[T any](w *World, e Entity) bool {
	return ecs.RemoveComponent[T](w, e)
}

// ------
// system
// ------

// a system operates on the components of entitues
type System = ecs.System

// `GetSystemId` retrieves the unique identifier of the current running system
//
// the id is provided by the scheduler at system registration time (`xp.AddSystem`)
var GetSystemId = ecs.GetSystemId

// ---------
// resources
// ---------

// add a resource to the world
// it is recommended to add the resource by value and not by reference
// the resource will be stored internally as value
// as it uses go generics, there can only be one value per type (ecs global singletons)
func AddResource[T any](w *World, r T) {
	ecs.AddResource(w, r)
}

var AddResourceByType = ecs.AddResourceByType

// the type `T` of the function cannot be a reference to a type
// if it is, it will return false
// the result will always be a reference to the resource so that
// the user can mutate its value
func GetResource[T any](w *World) (*T, bool) {
	return ecs.GetResource[T](w)
}

// the type `T` of the function cannot be a reference to a type
// if it is, it will return false
func RemoveResorce[T any](w *World) bool {
	return ecs.RemoveResource[T](w)
}

// -----
// query
// -----

type Filter = ecs.Filter

func With[T any]() Filter {
	return ecs.With[T]()
}

func Without[T any]() Filter {
	return ecs.Without[T]()
}

var Or = ecs.Or

func Query1[A any](w *World, filters ...Filter) *ecs.Query1[A] {
	return ecs.NewQuery1[A](w, filters...)
}

func Query2[A, B any](w *World, filters ...Filter) *ecs.Query2[A, B] {
	return ecs.NewQuery2[A, B](w, filters...)
}

func Query3[A, B, C any](w *World, filters ...Filter) *ecs.Query3[A, B, C] {
	return ecs.NewQuery3[A, B, C](w, filters...)
}

func Query4[A, B, C, D any](w *World, filters ...Filter) *ecs.Query4[A, B, C, D] {
	return ecs.NewQuery4[A, B, C, D](w, filters...)
}

// =====
// event
// =====

// description
func AddEvent[T any](w *World, ev T) {
	event.AddEvent(w, ev)
}

// description
func GetEvents[T any](w *World) ([]T, bool) {
	return event.GetEvents[T](w)
}

// ====
// game
// ====

// ----
// game
// ----
type Game = game.Game

// ------
// plugin
// ------

// ========
// schedule
// ========

// ---------
// condition
// ---------

// a condition function is defined as
//
// func(*World) bool
//
// it takes the ecs world as a parameter and returns a boolean indicating if the system
// should or should not run
type ConditionFn = schedule.ConditionFn

// description
var Once = schedule.Once

// description
var OnceWhen = schedule.OnceWhen

// ---------
// scheduler
// ---------

// the `Scheduler` is the main struct to run
//
// methods:
// - `RunIf(condition ConditionFn)`: add a condition for the execution of a system
type Scheduler = schedule.Scheduler

// adds a system to an schedule
var AddSystem = schedule.AddSystem

// -----
// stage
// -----

// description
var PreStartup = schedule.PreStartup

// description
var Startup = schedule.Startup

// description
var PostStartup = schedule.PostStartup

// description
var First = schedule.First

// description
var PreUpdate = schedule.PreUpdate

// description
var FixedFirst = schedule.FixedFirst

// description
var FixedPreUpdate = schedule.FixedPreUpdate

// description
var FixedUpdate = schedule.FixedUpdate

// description
var FixedPostUpdate = schedule.FixedPostUpdate

// description
var FixedLast = schedule.FixedLast

// description
var Update = schedule.Update

// description
var PostUpdate = schedule.PostUpdate

// description
var Last = schedule.Last

// description
var PreDraw = schedule.PreDraw

// description
var Draw = schedule.Draw

// description
var PostDraw = schedule.PostDraw

// #############
// CORE PACKAGES
// #############

// =====
// state
// =====

// ---------
// condition
// ---------

// ldhn
func InState[T comparable](s T) ConditionFn {
	return state.InState(s)
}

// -----
// stage
// -----

// description
func OnExit[T comparable](s T) schedule.StageFn {
	return state.OnExit(s)
}

// description
func OnTransition[T comparable](from, to T) schedule.StageFn {
	return state.OnTransition(from, to)
}

// description
func OnEnter[T comparable](s T) schedule.StageFn {
	return state.OnEnter(s)
}

// -----
// state
// -----

// description
func AddStateMachine[T comparable](sch *Scheduler, initial T) {
	state.AddStateMachine(sch, initial)
}

// description
func GetState[T comparable](w *World) (T, bool) {
	return state.GetState[T](w)
}

// description
func SetNextState[T comparable](w *World, s T) bool {
	return state.SetNextState(w, s)
}
